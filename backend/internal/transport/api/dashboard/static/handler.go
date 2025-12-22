package static

import (
	"net/http"
	"strings"
)

func RegisterFrontend(mux *http.ServeMux) {
	fileSystem := getFrontendFileSystem()
	fileServer := http.FileServer(fileSystem)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		f, err := fileSystem.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			// Try to serve index.html for the root
			r.URL.Path = "/"
		} else {
			_ = f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}
