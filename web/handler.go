package web

import (
	"embed"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/inventory/types"
)

// embed the dist folder
//
//go:embed static/*
var Static embed.FS

func NewInventoryHandler(storage storage.Storage) InventoryHandler {
	// Replace this in-memory function with a call to a database.
	reportsGetter := func() (reports []types.Report, err error) {
		return storage.GetAllReports(), nil
	}
	return InventoryHandler{
		GetReports: reportsGetter,
		Log:        log.Default(),
	}
}

type InventoryHandler struct {
	Log        *log.Logger
	GetReports func() ([]types.Report, error)
}

func (ph InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ps, err := ph.GetReports()
	// check the host query string
	host := r.URL.Query().Get("host")
	if host != "" {
		ph.Log.Printf("Filtering reports by host: %s", host)
		// filter the reports by host
		var filteredReports []types.Report
		for _, report := range ps {
			if report.Host.HostName == host {
				filteredReports = append(filteredReports, report)
			}
		}
		ps = filteredReports
	}
	container := r.URL.Query().Get("container")
	if container != "" {
		ph.Log.Printf("Filtering reports by container: %s", container)
		var filteredReports []types.Report

		// filter the reports by container
		for _, report := range ps {
			for _, c := range report.Containers {
				if c.ContainerID == container {
					// clear the container slice
					report.Containers = []types.Container{c}
					filteredReports = append(filteredReports, report)

				}
			}
		}
		ps = filteredReports

	}
	if err != nil {
		ph.Log.Printf("failed to get reports: %v", err)
		http.Error(w, "failed to retrieve reports", http.StatusInternalServerError)
		return
	}
	templ.Handler(reports(ps)).ServeHTTP(w, r)
}
