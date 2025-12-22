//go:build prod

package gewish

import "embed"

//go:embed all:frontend/build
var FrontendAssets embed.FS
