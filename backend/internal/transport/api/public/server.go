package public

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/LuukBlankenstijn/gewish/internal/app"
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
		Addr:      "localhost:8081",
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
		log.Default().Printf("Warning: no redirect found for \"%s%s\"", host, r.URL.String())
		http.NotFound(w, r)
		return
	}
	log.Default().Printf("Info: redirecting \"%s%s\" to \"%s\"", host, r.URL.String(), location)
	http.Redirect(w, r, location, http.StatusFound)
}
