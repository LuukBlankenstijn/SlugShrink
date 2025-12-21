package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/LuukBlankenstijn/gewish/internal/app"
	"github.com/LuukBlankenstijn/gewish/internal/authconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var permissionMap = map[apiv1.UserPermission]authconfig.UserPermission{
	apiv1.UserPermission_PERMISSION_USER:        authconfig.User,
	apiv1.UserPermission_PERMISSION_SUPERUSER:   authconfig.SuperUser,
	apiv1.UserPermission_PERMISSION_ADMIN:       authconfig.Admin,
	apiv1.UserPermission_PERMISSION_UNSPECIFIED: authconfig.User,
}

func AuthorizationInterceptor(authConfigs app.AuthConfigs) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			permissions, err := getAllowedPermissions(req.Spec().Procedure)
			if err == nil && len(permissions) == 0 {
				return next(ctx, req)
			}
			authConfig, err := authConfigs.Get(ctx)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			authState, isAuthenticated := authConfig.GetAuthState(ctx)
			if !isAuthenticated {
				slog.Error("tried to check permissions while skip_authentication is true", slog.String("procedure", req.Spec().Procedure))
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			for _, permission := range permissions {
				if (permission) == authState.UserPermission {
					return next(ctx, req)
				}
			}
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("not enough permissions"))
		}
	})
}

// decides based on a proto options if some request is allowed to skip authentication
func getAllowedPermissions(procedure string) (permissions []authconfig.UserPermission, err error) {
	if !strings.HasPrefix(procedure, "/") {
		return []authconfig.UserPermission{}, fmt.Errorf("bad procedure: %q", procedure)
	}
	p := strings.TrimPrefix(procedure, "/")
	i := strings.LastIndexByte(p, '/')
	if i < 0 {
		return []authconfig.UserPermission{}, fmt.Errorf("bad procedure: %q", procedure)
	}

	svcName := protoreflect.FullName(p[:i])
	mthName := protoreflect.Name(p[i+1:])

	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(svcName)
	if err != nil {
		return []authconfig.UserPermission{}, err
	}

	svc := desc.(protoreflect.ServiceDescriptor)
	md := svc.Methods().ByName(mthName)
	if md == nil {
		return []authconfig.UserPermission{}, fmt.Errorf("method not found: %s/%s", svcName, mthName)
	}

	opts := md.Options().(*descriptorpb.MethodOptions)

	v := proto.GetExtension(opts, apiv1.E_Permissions)

	b, ok2 := v.([]apiv1.UserPermission)
	if !ok2 {
		return []authconfig.UserPermission{}, fmt.Errorf("skip_auth has unexpected type")
	}
	permissions = []authconfig.UserPermission{}
	for _, p := range b {
		permissions = append(permissions, permissionMap[p])
	}
	return permissions, nil
}
