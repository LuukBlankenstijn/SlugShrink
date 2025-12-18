package app

import (
	"context"

	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"github.com/google/uuid"
)

type Domains struct {
	repo *gormrepo.DomainsRepo
}

func NewDomains(repo *gormrepo.DomainsRepo) Domains {
	return Domains{
		repo,
	}
}

func (a *Domains) Get(ctx context.Context, id uuid.UUID) (*apiv1.Domain, error) {
	m, err := a.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &apiv1.Domain{
		Id:     m.ID.String(),
		Name:   m.Name,
		Domain: m.Domain,
	}, nil
}

func (a *Domains) GetMany(ctx context.Context, pageData *apiv1.DomainsRequest) (*apiv1.DomainsResponse, error) {
	if pageData.Page == 0 {
		pageData.Page = 1
	}
	if pageData.Pagesize == 0 {
		pageData.Pagesize = 10
	}
	domains, totalRecords, err := a.repo.GetMany(ctx, int(pageData.Page), int(pageData.Pagesize))
	data := []*apiv1.Domain{}
	for _, domain := range domains {
		data = append(data, &apiv1.Domain{
			Id:     domain.ID.String(),
			Name:   domain.Name,
			Domain: domain.Domain,
		})
	}

	return &apiv1.DomainsResponse{
		Data:  data,
		Total: totalRecords,
	}, err
}

func (a *Domains) Delete(ctx context.Context, id uuid.UUID) error {
	return a.repo.Delete(ctx, id)
}

func (a *Domains) Create(ctx context.Context, domain *apiv1.CreateDomainRequest) (*apiv1.Domain, error) {
	newDomain, err := a.repo.Create(ctx, &models.DomainModel{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      domain.Name,
		Domain:    domain.Domain,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.Domain{
		Id:     newDomain.ID.String(),
		Name:   newDomain.Name,
		Domain: newDomain.Domain,
	}, nil
}

func (a *Domains) Put(ctx context.Context, domain *apiv1.Domain) (*apiv1.Domain, error) {
	id, err := uuid.Parse(domain.Id)
	if err != nil {
		return nil, err
	}

	updatedModel, err := a.repo.Save(ctx, models.DomainModel{
		BaseModel: models.BaseModel{ID: id},
		Name:      domain.Name,
		Domain:    domain.Domain,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.Domain{
		Id:     updatedModel.ID.String(),
		Name:   updatedModel.Name,
		Domain: updatedModel.Domain,
	}, nil
}
