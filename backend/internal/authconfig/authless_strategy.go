package authconfig

import (
	"context"

	"connectrpc.com/connect"
)

const authlessContextKey contextKey = "authless"

type AuthlessStrategy struct{}

func (AuthlessStrategy) Type() AuthStrategyType { return AuthStrategyAuthless }

func (AuthlessStrategy) Authenticate(ctx context.Context, _ connect.AnyRequest) (context.Context, error) {
	ctx = context.WithValue(ctx, authlessContextKey, true)
	return ctx, nil
}

func (AuthlessStrategy) IsAuthenticated(ctx context.Context) bool {
	value, ok := ctx.Value(authlessContextKey).(bool)
	return value && ok
}
