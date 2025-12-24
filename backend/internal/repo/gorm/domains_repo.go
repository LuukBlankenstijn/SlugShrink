package gormrepo

import (
	models "github.com/LuukBlankenstijn/slugshrink/internal/repo/gorm/models"
	"gorm.io/gorm"
)

type DomainsRepo struct {
	// This embeds all methods from BaseRepo into DomainsRepo
	BaseRepo[models.DomainModel]
}

func NewDomainsRepo(db *gorm.DB) *DomainsRepo {
	return &DomainsRepo{
		BaseRepo: NewRepo[models.DomainModel](db),
	}
}
