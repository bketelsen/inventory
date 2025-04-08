// Package web provides the HTTP handlers for the inventory service.
// It uses the templ package for rendering HTML templates.
// It also serves static files from the static directory.
package web

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/bketelsen/inventory"
)

// embed the dist folder
//
//go:embed static/*
var Static embed.FS

// NewInventoryHandler creates a new InventoryHandler
// It takes a storage interface as an argument
// and returns an InventoryHandler that uses the storage
func NewInventoryHandler(storage inventory.Storage) InventoryHandler {
	reportsGetter := func() (reports []inventory.Report, err error) {
		return storage.GetAllReports(), nil
	}
	return InventoryHandler{
		GetReports: reportsGetter,
	}
}

// InventoryHandler is the HTTP handler for the inventory service
type InventoryHandler struct {
	GetReports func() ([]inventory.Report, error)
}

// ServeHTTP implements the http.Handler interface
func (ph InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ps, err := ph.GetReports()
	if err != nil {
		slog.Error("failed to get reports", "error", err)
		http.Error(w, "failed to retrieve reports", http.StatusInternalServerError)
		return
	}
	// check the host query string
	host := r.URL.Query().Get("host")
	if host != "" {
		slog.Debug("Filtering reports", "host", host)
		// filter the reports by host
		var filteredReports []inventory.Report
		for _, report := range ps {
			if report.Host.HostName == host {
				filteredReports = append(filteredReports, report)
			}
		}
		ps = filteredReports
	}
	container := r.URL.Query().Get("container")
	if container != "" {
		slog.Debug("Filtering reports", "container", container)
		var filteredReports []inventory.Report

		// filter the reports by container
		for _, report := range ps {
			for _, c := range report.Containers {
				if c.ContainerID == container {
					// clear the container slice
					report.Containers = []inventory.Container{c}
					filteredReports = append(filteredReports, report)

				}
			}
		}
		ps = filteredReports
	}

	templ.Handler(reports(ps)).ServeHTTP(w, r)
}
