package app

import (
	"context"

	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"github.com/google/uuid"
)

type Redirects struct {
	repo *gormrepo.RedirectRepo
}

func NewRedirects(repo *gormrepo.RedirectRepo) Redirects {
	return Redirects{
		repo,
	}
}

func (a *Redirects) Get(ctx context.Context, id uuid.UUID) (*apiv1.Redirect, error) {
	m, err := a.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &apiv1.Redirect{
		DomainId:  m.DomainId.String(),
		Path:      m.Path,
		TargetUrl: m.TargetUrl,
		Active:    m.Active,
	}, nil
}

func (a *Redirects) Delete(ctx context.Context, id uuid.UUID) error {
	return a.repo.Delete(ctx, id)
}

func (a *Redirects) Create(ctx context.Context, redirect *apiv1.CreateRedirectRequest) (*apiv1.Redirect, error) {
	domainId, err := uuid.Parse(redirect.DomainId)
	if err != nil {
		return nil, err
	}
	newRedirect, err := a.repo.Create(ctx, &models.RedirectModel{
		BaseModel: models.BaseModel{ID: uuid.New()},
		DomainId:  domainId,
		Path:      redirect.Path,
		TargetUrl: redirect.TargetUrl,
		Active:    redirect.Active,
	})
	return &apiv1.Redirect{
		DomainId:  redirect.DomainId,
		Path:      redirect.Path,
		TargetUrl: redirect.TargetUrl,
		Active:    redirect.Active,
		Id:        newRedirect.ID.String(),
	}, err
}

func (a *Redirects) Put(ctx context.Context, redirect *apiv1.Redirect) (*apiv1.Redirect, error) {
	id, err := uuid.Parse(redirect.Id)
	if err != nil {
		return nil, err
	}
	domainId, err := uuid.Parse(redirect.DomainId)
	if err != nil {
		return nil, err
	}
	updatedModel, err := a.repo.Save(ctx, models.RedirectModel{
		BaseModel: models.BaseModel{ID: id},
		DomainId:  domainId,
		Path:      redirect.Path,
		TargetUrl: redirect.TargetUrl,
		Active:    redirect.Active,
	})
	return &apiv1.Redirect{
		Id:        updatedModel.ID.String(),
		DomainId:  updatedModel.DomainId.String(),
		Path:      updatedModel.Path,
		TargetUrl: updatedModel.TargetUrl,
		Active:    updatedModel.Active,
	}, err
}
