package authconfig

import (
	"context"

	"connectrpc.com/connect"
)

type AuthStrategy interface {
	Type() AuthStrategyType
	Authenticate(context context.Context, req connect.AnyRequest) (context.Context, error)
}
