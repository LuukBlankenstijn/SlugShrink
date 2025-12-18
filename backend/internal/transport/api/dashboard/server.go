package dashboard

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/LuukBlankenstijn/gewish/gen/api/v1/apiv1connect"
	"github.com/LuukBlankenstijn/gewish/internal/app"
)

type DashboardApi struct {
	redirects app.Redirects
	domains   app.Domains
}

func NewDashboardApi(redirects app.Redirects, domains app.Domains) DashboardApi {
	return DashboardApi{
		redirects: redirects,
		domains:   domains,
	}
}

func (a *DashboardApi) Run() error {
	mux := http.NewServeMux()
	path, handler := apiv1connect.NewApiServiceHandler(
		a,
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	mux.Handle(path, handler)

	reflector := grpcreflect.NewStaticReflector(
		"api.v1.ApiService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}

	return s.ListenAndServe()
}
