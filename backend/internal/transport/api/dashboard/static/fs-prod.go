//go:build prod

package static

import (
	"io/fs"
	"net/http"

	"github.com/LuukBlankenstijn/gewish"
)

func getFrontendFileSystem() http.FileSystem {
	distFS, _ := fs.Sub(gewish.FrontendAssets, "frontend/build")
	return http.FS(distFS)
}
