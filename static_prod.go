//go:build !dev
// +build !dev

package inventory

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	hashFS "github.com/benbjohnson/hashfs"
)

//go:embed web/static/*
var staticFS embed.FS
var staticRootFS, _ = fs.Sub(staticFS, "web/static")

func Static(logger *slog.Logger) http.Handler {
	logger.Debug("static assets are embedded")
	return hashFS.FileServer(staticRootFS)
}
