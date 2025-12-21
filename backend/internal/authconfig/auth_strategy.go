package authconfig

import (
	"context"

	"connectrpc.com/connect"
)

type AuthStrategy interface {
	Type() AuthStrategyType
	Authenticate(context context.Context, req connect.AnyRequest) (context.Context, error)
	GetAuthState(ctx context.Context) (*authContext, bool)
}
