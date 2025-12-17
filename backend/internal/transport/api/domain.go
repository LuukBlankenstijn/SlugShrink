package api

import (
	"context"

	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ApiServer) GetDomain(
	ctx context.Context,
	req *apiv1.DomainRequest,
) (*apiv1.Domain, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return s.domains.Get(ctx, id)
}

func (s *ApiServer) DeleteDomain(
	ctx context.Context,
	req *apiv1.DomainRequest,
) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), s.domains.Delete(ctx, id)
}

func (s *ApiServer) PutDomain(
	ctx context.Context,
	req *apiv1.Domain,
) (*apiv1.Domain, error) {
	return s.domains.Put(ctx, req)
}

func (s *ApiServer) CreateDomain(
	ctx context.Context,
	req *apiv1.CreateDomainRequest,
) (*apiv1.Domain, error) {
	return s.domains.Create(ctx, req)
}
