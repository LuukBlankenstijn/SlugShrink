package authconfig

import (
	"context"

	"connectrpc.com/connect"
)

type AuthlessStrategy struct {
	baseStrategy
}

func (AuthlessStrategy) Type() AuthStrategyType { return AuthStrategyAuthless }

func (s AuthlessStrategy) Authenticate(ctx context.Context, _ connect.AnyRequest) (context.Context, error) {
	ctx = s.setContext(ctx, Admin, nil)
	return ctx, nil
}
