// Package webui embeds the built web UI. `make build` copies the SvelteKit
// output into dist/ before compiling; a plain `go build` embeds only .gitkeep.
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// FS returns the web UI files rooted at the build output directory.
func FS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
