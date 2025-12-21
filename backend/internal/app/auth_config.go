package app

import (
	"context"

	"github.com/LuukBlankenstijn/gewish/internal/authconfig"
	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
)

type AuthConfigs struct {
	repo *gormrepo.AuthConfigRepo
}

func NewAuthConfigs(repo *gormrepo.AuthConfigRepo) AuthConfigs {
	return AuthConfigs{
		repo,
	}
}

func (a *AuthConfigs) Get(ctx context.Context) (authconfig.AuthStrategy, error) {
	return a.repo.Get(ctx)
}

func (a *AuthConfigs) Set(ctx context.Context, strategy authconfig.AuthStrategy) error {
	return a.repo.Set(ctx, strategy)
}

func (a *AuthConfigs) Login(ctx context.Context, password string) (string, error) {
	config, err := a.Get(ctx)
	if err != nil {
		return "", err
	}
	switch stategy := config.(type) {
	case *authconfig.BasicStrategy:
		return stategy.Login(password)
	default:
		return "", nil

	}
}
