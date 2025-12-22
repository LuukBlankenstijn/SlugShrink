//go:build !prod

package static

import (
	"embed"
	"net/http"
)

var ui embed.FS

func getFrontendFileSystem() http.FileSystem {
	return http.FS(ui)
}
