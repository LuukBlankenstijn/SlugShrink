package dashboard

import (
	"context"

	apiv1 "github.com/LuukBlankenstijn/slugshrink/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *DashboardApi) GetRedirect(
	ctx context.Context,
	req *apiv1.RedirectRequest,
) (*apiv1.Redirect, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return s.redirects.Get(ctx, id)
}

func (s *DashboardApi) GetRedirects(
	ctx context.Context,
	req *apiv1.RedirectsRequest,
) (*apiv1.RedirectsResponse, error) {
	return s.redirects.GetMany(ctx, req)
}

func (s *DashboardApi) DeleteRedirect(
	ctx context.Context,
	req *apiv1.RedirectRequest,
) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	return new(emptypb.Empty), s.redirects.Delete(ctx, id)
}

func (s *DashboardApi) PutRedirect(
	ctx context.Context,
	req *apiv1.Redirect,
) (*apiv1.Redirect, error) {
	return s.redirects.Put(ctx, req)
}

func (s *DashboardApi) CreateRedirect(
	ctx context.Context,
	req *apiv1.CreateRedirectRequest,
) (*apiv1.Redirect, error) {
	return s.redirects.Create(ctx, req)
}
