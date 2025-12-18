package types

import (
	"context"

	"connectrpc.com/connect"
)

type contextKey string

type AuthStrategy interface {
	Type() AuthStrategyType
	Authenticate(context context.Context, req connect.AnyRequest) (context.Context, error)
	IsAuthenticated(context context.Context) bool
}
