package gormrepo

import (
	"context"
	"path"
	"strings"

	models "github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *RedirectRepo) FindRedirectForHostAndPath(context context.Context, host, rawpath string) (models.RedirectModel, error) {
	cleanPath := path.Clean(rawpath)
	searchPaths := []string{cleanPath, strings.TrimSuffix(cleanPath, "/") + "/"}

	redirect, err := r.repo.
		Joins(clause.Has("Domain"), func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Where(map[string]any{"domain": host})
			return nil
		}).
		Where("path IN ?", searchPaths).
		First(context)

	return redirect, err
}
