package assets

import "embed"

//go:embed goshrt.png style.css
var PublicAssets embed.FS

//go:embed landingpage.tmpl
var InternalAssets embed.FS
