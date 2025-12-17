package api

import (
	"context"

	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ApiServer) GetRedirect(
	ctx context.Context,
	req *apiv1.RedirectRequest,
) (*apiv1.Redirect, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return s.redirects.Get(ctx, id)
}

func (s *ApiServer) DeleteRedirect(
	ctx context.Context,
	req *apiv1.RedirectRequest,
) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return new(emptypb.Empty), s.redirects.Delete(ctx, id)
}

func (s *ApiServer) PutRedirect(
	ctx context.Context,
	req *apiv1.Redirect,
) (*apiv1.Redirect, error) {
	return s.redirects.Put(ctx, req)
}

func (s *ApiServer) CreateRedirect(
	ctx context.Context,
	req *apiv1.CreateRedirectRequest,
) (*apiv1.Redirect, error) {
	return s.redirects.Create(ctx, req)
}
