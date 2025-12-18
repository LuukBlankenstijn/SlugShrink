package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepo[T any] struct {
	repo gorm.Interface[T]
}

func NewRepo[T any](db *gorm.DB) BaseRepo[T] {
	return BaseRepo[T]{
		gorm.G[T](db),
	}
}

func (r *BaseRepo[T]) Get(context context.Context, id uuid.UUID) (*T, error) {
	m, err := r.repo.Where("id = ?", id.String()).First(context)
	return &m, err
}

func (r *BaseRepo[T]) GetMany(context context.Context, page, pagesize int) ([]T, int64, error) {
	models, err := r.repo.Limit(pagesize).Offset((page - 1) * pagesize).Find(context)
	if err != nil {
		return []T{}, 0, err
	}
	count, err := r.repo.Count(context, "id")

	return models, count, err
}

func (r *BaseRepo[T]) Create(context context.Context, newModel *T) (*T, error) {
	err := r.repo.Create(context, newModel)

	return newModel, err
}

func (r *BaseRepo[T]) Save(context context.Context, model T) (*T, error) {
	_, err := r.repo.Updates(context, model)
	return &model, err
}

func (r *BaseRepo[T]) Delete(context context.Context, id uuid.UUID) error {
	_, err := r.repo.Where("id = ?", id.String()).Delete(context)
	return err
}
