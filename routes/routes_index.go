package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bketelsen/inventory"
	"github.com/bketelsen/inventory/web/components"
	"github.com/bketelsen/inventory/web/pages"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/go-chi/chi/v5"
)

func setupIndexRoute(router chi.Router, storage inventory.Storage) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Index getting reports")
		reports := storage.GetAllReports()
		host := r.URL.Query().Get("host")
		if host != "" {
			slog.Info("Filtering reports", "host", host)
			newReports := make([]inventory.Report, 0)

			// filter the reports by host
			for _, report := range reports {
				if report.Host.HostName == host {
					newReports = append(newReports, report)
				}
			}
			reports = newReports
		}
		platform := r.URL.Query().Get("platform")
		if platform != "" {
			newReports := make([]inventory.Report, 0)
			slog.Info("Filtering reports", "platform", platform)

			// filter the reports by container
			for _, report := range reports {
				if report.HasPlatform(platform) {
					newReports = append(newReports, report)
				}
			}
			reports = newReports
		}

		slog.Info("Dashboard Initial", "reports", len(reports))
		_ = pages.DashboardInitial(reports, platform, host).Render(r.Context(), w)
	})

	router.Get("/dashboard/data", func(w http.ResponseWriter, r *http.Request) {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		sse := datastar.NewSSE(w, r)

		signals := &PageSignals{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reports := storage.GetAllReports()
		host := signals.Host
		if host != "" {
			slog.Info("Filtering reports", "host", host)
			newReports := make([]inventory.Report, 0)

			// filter the reports by host
			for _, report := range reports {
				if report.Host.HostName == host {
					newReports = append(newReports, report)
				}
			}
			reports = newReports
		}
		platform := signals.Platform
		if platform != "" {
			newReports := make([]inventory.Report, 0)
			slog.Info("Filtering reports", "platform", platform)

			// filter the reports by container
			for _, report := range reports {
				if report.HasPlatform(platform) {
					newReports = append(newReports, report)
				}
			}
			reports = newReports
		}
		_ = sse.MergeFragmentTempl(pages.Filters(reports))
		_ = sse.MergeFragmentTempl(components.Dashboard(reports))

		for {
			select {
			case <-r.Context().Done():
				slog.Info("client disconnected")
				return

			case <-ticker.C:
				signals := &PageSignals{}
				if err := datastar.ReadSignals(r, signals); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				slog.Info("Dashboard Data - New Signal", "signals", signals)
				slog.Info("Dashboard Data poller getting reports")

				reports := storage.GetAllReports()
				host := signals.Host
				if host != "" {
					slog.Info("Filtering reports", "host", host)
					newReports := make([]inventory.Report, 0)

					// filter the reports by host
					for _, report := range reports {
						if report.Host.HostName == host {
							newReports = append(newReports, report)
						}
					}
					reports = newReports
				}
				platform := signals.Platform
				if platform != "" {
					newReports := make([]inventory.Report, 0)
					slog.Info("Filtering reports", "platform", platform)

					// filter the reports by container
					for _, report := range reports {
						if report.HasPlatform(platform) {
							newReports = append(newReports, report)
						}
					}
					reports = newReports
				}
				_ = sse.MergeFragmentTempl(pages.Filters(reports))
				_ = sse.MergeFragmentTempl(components.Dashboard(reports))
			}
		}
	})
}

type PageSignals struct {
	Host     string `json:"host"`
	Platform string `json:"platform"`
}
