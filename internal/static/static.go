package static

import (
	"embed"
	"io/fs"
)

//go:embed css/* js/* swagger.html openapi.yaml
var staticFS embed.FS

func FS() fs.FS {
	return staticFS
}
