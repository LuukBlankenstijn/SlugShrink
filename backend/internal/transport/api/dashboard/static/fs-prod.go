//go:build prod

package static

import (
	"io/fs"
	"net/http"

	"github.com/LuukBlankenstijn/slugshrink"
)

func getFrontendFileSystem() http.FileSystem {
	distFS, _ := fs.Sub(slugshrink.FrontendAssets, "frontend/build")
	return http.FS(distFS)
}
