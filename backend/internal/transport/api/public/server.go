package public

import (
	"errors"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/LuukBlankenstijn/slugshrink/internal/app"
	"github.com/LuukBlankenstijn/slugshrink/internal/logging"
	"gorm.io/gorm"
)

type PublicApi struct {
	redirects app.Redirects
	domains   app.Domains
}

func NewPublicApi(redirects app.Redirects, domains app.Domains) PublicApi {
	return PublicApi{
		redirects: redirects,
		domains:   domains,
	}
}

func (a *PublicApi) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.redirectHandler)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "0.0.0.0:8081",
		Handler:   mux,
		Protocols: p,
	}

	return s.ListenAndServe()
}

func (p *PublicApi) redirectHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	path := path.Clean(r.URL.Path)
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}
	location, err := p.redirects.FindLocationByHostAndPath(r.Context(), host, path)
	if err != nil {
		level := slog.LevelError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			level = slog.LevelWarn
		}
		logging.Logger().Log(r.Context(), level, "failed to resolve redirect", slog.String("host", host), slog.String("path", r.URL.String()), slog.Any("error", err))
		http.NotFound(w, r)
		return
	}
	logging.Logger().Info("redirecting request", slog.String("host", host), slog.String("path", r.URL.String()), slog.String("location", location))
	http.Redirect(w, r, location, http.StatusFound)
}
