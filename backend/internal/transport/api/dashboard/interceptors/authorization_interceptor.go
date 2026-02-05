package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	apiv1 "github.com/LuukBlankenstijn/slugshrink/gen/api/v1"
	"github.com/LuukBlankenstijn/slugshrink/internal/app"
	"github.com/LuukBlankenstijn/slugshrink/internal/authconfig"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var permissionMap = map[apiv1.UserPermission]authconfig.UserPermission{
	apiv1.UserPermission_PERMISSION_USER:      authconfig.User,
	apiv1.UserPermission_PERMISSION_SUPERUSER: authconfig.SuperUser,
	apiv1.UserPermission_PERMISSION_ADMIN:     authconfig.Admin,
}

type permissionRequirements struct {
	allowed        map[authconfig.UserPermission]struct{}
	requireCreator bool
}

func newPermissionRequirements(perms []apiv1.UserPermission) permissionRequirements {
	req := permissionRequirements{
		allowed:        make(map[authconfig.UserPermission]struct{}),
		requireCreator: false,
	}
	for _, p := range perms {
		if p == apiv1.UserPermission_PERMISSION_CREATOR {
			req.requireCreator = true
			continue
		}
		if mapped, ok := permissionMap[p]; ok {
			req.allowed[mapped] = struct{}{}
		}
	}
	return req
}

func (p permissionRequirements) allows(permission authconfig.UserPermission) bool {
	_, ok := p.allowed[permission]
	return ok
}

func AuthorizationInterceptor(authConfigs app.AuthConfigs, redirects app.Redirects) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Resolve permissions from proto options; if none are required, continue.
			permissions, err := getAllowedPermissions(req.Spec().Procedure)
			if err == nil && len(permissions) == 0 {
				return next(ctx, req)
			}
			reqPerms := newPermissionRequirements(permissions)
			// Tag creator requirement if the request carries a redirect ID.
			redirectID, hasRedirectID := extractRedirectID(req)
			if hasRedirectID {
				reqPerms.requireCreator = true
			}

			// Pull auth state from context (populated by authentication interceptor).
			authState, isAuthenticated := authconfig.GetAuthState(ctx)
			if !isAuthenticated {
				slog.Error("tried to check permissions while skip_authentication is true", slog.String("procedure", req.Spec().Procedure))
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			// Role-based allow first.
			if reqPerms.allows(authState.UserPermission) {
				return next(ctx, req)
			}
			// If role is insufficient but creator access is allowed, validate ownership.
			if reqPerms.requireCreator {
				if redirectID != nil && authState.UserId != nil {
					slog.Debug("test", slog.String("redirect", redirectID.String()), slog.String("user", *authState.UserId))
				}
				if redirectID != nil && authState.UserId != nil && redirects.IsOwner(ctx, *redirectID, *authState.UserId) {
					return next(ctx, req)
				}
			}
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("not enough permissions"))
		}
	})
}

func getAllowedPermissions(procedure string) (permissions []apiv1.UserPermission, err error) {
	if !strings.HasPrefix(procedure, "/") {
		return []apiv1.UserPermission{}, fmt.Errorf("bad procedure: %q", procedure)
	}
	p := strings.TrimPrefix(procedure, "/")
	i := strings.LastIndexByte(p, '/')
	if i < 0 {
		return []apiv1.UserPermission{}, fmt.Errorf("bad procedure: %q", procedure)
	}

	svcName := protoreflect.FullName(p[:i])
	mthName := protoreflect.Name(p[i+1:])

	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(svcName)
	if err != nil {
		return []apiv1.UserPermission{}, err
	}

	svc := desc.(protoreflect.ServiceDescriptor)
	md := svc.Methods().ByName(mthName)
	if md == nil {
		return []apiv1.UserPermission{}, fmt.Errorf("method not found: %s/%s", svcName, mthName)
	}

	opts := md.Options().(*descriptorpb.MethodOptions)

	v := proto.GetExtension(opts, apiv1.E_Permissions)

	b, ok2 := v.([]apiv1.UserPermission)
	if !ok2 {
		return []apiv1.UserPermission{}, fmt.Errorf("skip_auth has unexpected type")
	}
	return b, nil
}

func extractRedirectID(req connect.AnyRequest) (*uuid.UUID, bool) {
	// Walk fields on the request message and return a UUID for any field tagged with is_redirect_id.
	msg, ok := req.Any().(interface {
		ProtoReflect() protoreflect.Message
	})
	if !ok {
		return nil, false
	}

	m := msg.ProtoReflect()
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		opts, ok := fd.Options().(*descriptorpb.FieldOptions)
		if !ok {
			continue
		}
		raw := proto.GetExtension(opts, apiv1.E_IsRedirectId)
		isRedirect, ok := raw.(bool)
		if !ok || !isRedirect {
			continue
		}
		if !m.Has(fd) {
			continue
		}
		if fd.Kind() == protoreflect.StringKind && m.Get(fd).String() != "" {
			value := m.Get(fd).String()
			id, err := uuid.Parse(value)
			if err != nil {
				continue
			}
			return &id, true
		}
	}
	return nil, false
}
