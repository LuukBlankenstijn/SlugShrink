package gormrepo

import (
	"context"
	"errors"

	"github.com/LuukBlankenstijn/gewish/internal/authconfig"
	models "github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthConfigRepo struct {
	repo gorm.Interface[models.AuthConfigModel]
}

func NewAuthConfigRepo(db *gorm.DB) *AuthConfigRepo {
	return &AuthConfigRepo{
		repo: gorm.G[models.AuthConfigModel](db),
	}
}

func (r *AuthConfigRepo) Set(context context.Context, strategy authconfig.AuthStrategy) error {
	existing, err := r.repo.First(context)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.repo.Create(context, &models.AuthConfigModel{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
			Strategy: strategy,
		})
	}
	if err != nil {
		return err
	}

	_, err = r.repo.Where("id = ?", existing.ID.String()).Delete(context)
	if err != nil {
		return err
	}

	return r.repo.Create(context, &models.AuthConfigModel{
		BaseModel: models.BaseModel{
			ID: existing.ID,
		},
		Strategy: strategy,
	})
}

func (r *AuthConfigRepo) Get(context context.Context) (authconfig.AuthStrategy, error) {
	var config models.AuthConfigModel
	config, err := r.repo.First(context)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return authconfig.AuthlessStrategy{}, nil
	}
	if err != nil {
		return nil, err
	}
	return config.Strategy, nil
}
