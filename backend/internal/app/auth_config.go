package app

import (
	"context"

	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/types"
)

type AuthConfigs struct {
	repo *gormrepo.AuthConfigRepo
}

func NewAuthConfigs(repo *gormrepo.AuthConfigRepo) AuthConfigs {
	return AuthConfigs{
		repo,
	}
}

func (a *AuthConfigs) Get(ctx context.Context) (types.AuthStrategy, error) {
	return a.repo.Get(ctx)
}

func (a *AuthConfigs) Set(ctx context.Context, strategy types.AuthStrategy) error {
	return a.repo.Set(ctx, strategy)
}

func (a *AuthConfigs) Login(ctx context.Context, password string) (string, error) {
	config, err := a.Get(ctx)
	if err != nil {
		return "", err
	}
	switch stategy := config.(type) {
	case *types.BasicStrategy:
		return stategy.Login(password)
	default:
		return "", nil

	}
}
