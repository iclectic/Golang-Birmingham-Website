// Package web embeds the site's templates and static assets into the binary.
package web

import (
	"embed"
	"io/fs"
)

//go:embed templates static
var files embed.FS

// Templates holds the html/template sources.
var Templates = mustSub("templates")

// Static holds CSS, JavaScript and images served under /static/.
var Static = mustSub("static")

func mustSub(dir string) fs.FS {
	sub, err := fs.Sub(files, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
