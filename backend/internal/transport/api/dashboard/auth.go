package dashboard

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	apiv1 "github.com/LuukBlankenstijn/gewish/gen/api/v1"
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/types"
)

func (s *DashboardApi) GetAuthStatus(ctx context.Context, _ *emptypb.Empty) (*apiv1.AuthStatus, error) {
	config, err := s.authConfigs.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &apiv1.AuthStatus{
		Authenticated: config.IsAuthenticated(ctx),
	}, nil
}

func (s *DashboardApi) Login(ctx context.Context, login *apiv1.BasicAuthLogin) (*emptypb.Empty, error) {
	callInfo, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, errors.New("can't access headers: no CallInfo for handler context")
	}
	token, err := s.authConfigs.Login(ctx, login.Password)
	if err != nil {
		return nil, err
	}
	// if no token and no error, the configured auth strategy does not need a token cookie
	if token == "" {
		return &emptypb.Empty{}, nil
	}
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
		Path:     "/",
	}
	if v := cookie.String(); v != "" {
		callInfo.ResponseHeader().Add("Set-Cookie", v)
	}

	return &emptypb.Empty{}, nil
}

func (s *DashboardApi) SetAuthConfig(ctx context.Context, setAuth *apiv1.SetAuth) (*emptypb.Empty, error) {
	var config types.AuthStrategy
	switch apiConfig := setAuth.Config.(type) {
	case *apiv1.SetAuth_BasicAuth:
		config = types.NewBasicStrategy(apiConfig.BasicAuth.Password)
	case *apiv1.SetAuth_Authless:
		config = types.AuthlessStrategy{}
	default:
		config = types.AuthlessStrategy{}
	}
	return &emptypb.Empty{}, s.authConfigs.Set(ctx, config)
}

func (s *DashboardApi) GetAuthConfig(ctx context.Context, empty *emptypb.Empty) (*apiv1.GetAuth, error) {
	config, err := s.authConfigs.Get(ctx)
	if err != nil {
		return nil, err
	}
	switch config.Type() {
	case types.AuthStrategyAuthless:
		return &apiv1.GetAuth{
			Config: &apiv1.GetAuth_Authless{},
		}, nil
	case types.AuthStrategyBasic:
		return &apiv1.GetAuth{
			Config: &apiv1.GetAuth_BasicAuth{},
		}, nil
	default:
		return nil, errors.New("unknow auth type")
	}
}
