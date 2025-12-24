package gormrepo

import (
	"context"
	"path"
	"strings"

	models "github.com/LuukBlankenstijn/slugshrink/internal/repo/gorm/models"
	"github.com/google/uuid"
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

func (r *RedirectRepo) GetFullRedirects(context context.Context, page, pagesize int, search string) ([]models.RedirectModel, int64, error) {
	query := r.repo.Preload("Domain", nil)
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(path) LIKE ? OR LOWER(target_url) LIKE ?", like, like)
	}

	found, err := query.Limit(pagesize).Offset((page - 1) * pagesize).Find(context)
	if err != nil {
		return []models.RedirectModel{}, 0, err
	}

	var count int64
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		count, err = r.repo.Where("LOWER(path) LIKE ? OR LOWER(target_url) LIKE ?", like, like).Count(context, "id")
	} else {
		count, err = r.repo.Count(context, "id")
	}

	return found, count, err
}

func (r *RedirectRepo) CheckIsOwner(context context.Context, user string, id uuid.UUID) bool {
	found, err := r.repo.Where("id = ?", id.String()).First(context)
	return err == nil && found.Creator != nil && *found.Creator == user
}
