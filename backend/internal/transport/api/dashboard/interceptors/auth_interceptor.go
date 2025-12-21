package interceptors

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/LuukBlankenstijn/gewish/internal/app"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func AuthInterceptor(authConfigs app.AuthConfigs) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			skip, _, _ := skipAuth(req.Spec().Procedure)
			authConfig, err := authConfigs.Get(ctx)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			ctx, err = authConfig.Authenticate(ctx, req)
			if err != nil && !skip {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			return next(ctx, req)
		}

	})
}

// decides based on a proto options if some request is allowed to skip authentication
func skipAuth(procedure string) (skip bool, ok bool, err error) {
	if !strings.HasPrefix(procedure, "/") {
		return false, false, fmt.Errorf("bad procedure: %q", procedure)
	}
	p := strings.TrimPrefix(procedure, "/")
	i := strings.LastIndexByte(p, '/')
	if i < 0 {
		return false, false, fmt.Errorf("bad procedure: %q", procedure)
	}

	svcName := protoreflect.FullName(p[:i])
	mthName := protoreflect.Name(p[i+1:])

	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(svcName)
	if err != nil {
		return false, false, err
	}

	svc := desc.(protoreflect.ServiceDescriptor)
	md := svc.Methods().ByName(mthName)
	if md == nil {
		return false, false, fmt.Errorf("method not found: %s/%s", svcName, mthName)
	}

	opts := md.Options().(*descriptorpb.MethodOptions)

	v := proto.GetExtension(opts, apiv1.E_SkipAuth)

	b, ok2 := v.(bool)
	if !ok2 {
		return false, false, fmt.Errorf("skip_auth has unexpected type")
	}
	return b, true, nil
}
