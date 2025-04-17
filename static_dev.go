//go:build dev
// +build dev

package inventory

import (
	"log/slog"
	"net/http"
	"os"
)

func Static(logger *slog.Logger) http.Handler {
	logger.Info("static assets are being served from web/static/")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.FileServerFS(os.DirFS("web/static")).ServeHTTP(w, r)
	})
}
