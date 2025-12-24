package app

import (
	"context"
	"errors"
	"strings"

	apiv1 "github.com/LuukBlankenstijn/slugshrink/gen/api/v1"
	"github.com/LuukBlankenstijn/slugshrink/internal/repo/gorm"
	"github.com/LuukBlankenstijn/slugshrink/internal/repo/gorm/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		Creator:   m.Creator,
	}, nil
}

func (a *Redirects) GetMany(ctx context.Context, pageData *apiv1.RedirectsRequest) (*apiv1.RedirectsResponse, error) {
	if pageData.Page == 0 {
		pageData.Page = 1
	}
	if pageData.Pagesize == 0 {
		pageData.Pagesize = 10
	}
	search := strings.TrimSpace(pageData.Search)
	redirects, total, err := a.repo.GetFullRedirects(ctx, int(pageData.Page), int(pageData.Pagesize), search)
	data := []*apiv1.FullRedirect{}

	for _, redirect := range redirects {
		data = append(data, &apiv1.FullRedirect{
			Id:        redirect.ID.String(),
			DomainId:  redirect.DomainId.String(),
			Domain:    redirect.Domain.Domain,
			Path:      redirect.Path,
			TargetUrl: redirect.TargetUrl,
			Active:    redirect.Active,
			Creator:   redirect.Creator,
		})
	}

	return &apiv1.RedirectsResponse{
		Data:  data,
		Total: int32(total),
	}, err
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
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, errors.New("this domain/path combination already exists")
	}
	if err != nil {
		return nil, err
	}
	return &apiv1.Redirect{
		DomainId:  redirect.DomainId,
		Path:      redirect.Path,
		TargetUrl: redirect.TargetUrl,
		Active:    redirect.Active,
		Id:        newRedirect.ID.String(),
		Creator:   newRedirect.Creator,
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
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, errors.New("this domain/path combination already exists")
	}
	if err != nil {
		return nil, err
	}
	return &apiv1.Redirect{
		Id:        updatedModel.ID.String(),
		DomainId:  updatedModel.DomainId.String(),
		Path:      updatedModel.Path,
		TargetUrl: updatedModel.TargetUrl,
		Active:    updatedModel.Active,
	}, err
}

func (a *Redirects) FindLocationByHostAndPath(ctx context.Context, host, path string) (string, error) {
	redirect, err := a.repo.FindRedirectForHostAndPath(ctx, host, path)
	return redirect.TargetUrl, err
}

func (a *Redirects) IsOwner(ctx context.Context, redirectID uuid.UUID, user string) bool {
	return a.repo.CheckIsOwner(ctx, user, redirectID)
}

func (a *Redirects) CheckIsCreator(ctx context.Context, user string, id uuid.UUID) bool {
	return a.repo.CheckIsOwner(ctx, user, id)
}
