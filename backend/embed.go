//go:build prod

package slugshrink

import "embed"

//go:embed all:frontend/build
var FrontendAssets embed.FS
