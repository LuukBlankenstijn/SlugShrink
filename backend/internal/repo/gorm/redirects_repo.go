package gormrepo

import (
	models "github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"gorm.io/gorm"
)

type RedirectRepo struct {
	// This embeds all methods from BaseRepo into DomainsRepo
	BaseRepo[models.RedirectModel]
}

func NewRedirectsRepo(db *gorm.DB) *RedirectRepo {
	return &RedirectRepo{
		BaseRepo: NewRepo[models.RedirectModel](db),
	}
}
