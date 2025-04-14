// Package service provides the inventory server implementation
// It implements the RPC methods for the inventory server
package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bketelsen/inventory"
)

// InventoryServer implements the inventory server interface
type InventoryServer struct {
	storage inventory.Storage
}

// NewInventoryServer creates a new server with the provided storage
func NewInventoryServer(storage inventory.Storage) *InventoryServer {
	return &InventoryServer{
		storage: storage,
	}
}

// Update is the method that the client will call to send a report
func (c *InventoryServer) Update(report *inventory.Report, reply *int) error {
	if report == nil {
		return inventory.ErrEmptyReport
	}
	report.Timestamp = time.Now()

	slog.Info("Received report", "host", report.Host.HostName, "ip", report.Host.IP)
	for _, l := range report.Listeners {
		slog.Debug("Listener", "host", report.Host.HostName, "address", l.ListenAddress, "pid", l.PID, "port", l.Port, "program", l.Program, "app", l.Program)
	}
	for _, ct := range report.Containers {
		slog.Debug("Container", "host", report.Host.HostName, "container", ct.ContainerID, "ip", ct.IP.String(), "hostname", ct.HostName, "platform", ct.Platform, "image", ct.Image)
	}
	for _, s := range report.Services {
		slog.Debug("Service", "host", report.Host.HostName, "name", s.Name, "port", s.Port, "protocol", s.Protocol)
	}

	// Store the report in persistent storage
	err := c.storage.StoreReport(*report)
	if err != nil {
		slog.Error("Failed to store report", "error", err)
		return fmt.Errorf("failed to store report: %w", err)
	}

	*reply = 0
	return nil
}

// Search is the method that the client will call to search for reports
// It takes a query string and returns a list of reports that match the query
func (c *InventoryServer) Search(query string, reply *[]inventory.Report) error {
	if query == "" {
		return inventory.ErrEmptyQuery
	}

	slog.Info("Received search query", "query", query)
	// Perform the search in persistent storage
	reports := c.storage.GetAllReports()
	if reports == nil {
		return inventory.ErrNoResults
	}
	for _, report := range reports {
		found := false
		rpt := inventory.Report{
			Host:       report.Host,
			Services:   []inventory.Service{},
			Containers: []inventory.Container{},
			Listeners:  []inventory.Listener{},
		}
		// iterate through reports and find matches on services and containers
		for _, service := range report.Services {
			lowerServiceName := strings.ToLower(service.Name)
			if strings.Contains(lowerServiceName, strings.ToLower(query)) {
				slog.Debug("Found service", "service", service.Name, "host", report.Host.HostName)
				rpt.Services = append(rpt.Services, service)
				found = true
			}
		}
		for _, container := range report.Containers {
			lowerContainerName := strings.ToLower(container.ContainerID)
			lowerImageName := strings.ToLower(container.Image)
			lowerHostName := strings.ToLower(container.HostName)
			if strings.Contains(lowerContainerName, strings.ToLower(query)) || strings.Contains(lowerImageName, strings.ToLower(query)) || strings.Contains(lowerHostName, strings.ToLower(query)) {
				slog.Debug("Found container", "container", container.ContainerID, "host", report.Host.HostName)
				rpt.Containers = append(rpt.Containers, container)
				found = true
			}
		}
		for _, listener := range report.Listeners {
			lowerListenerName := strings.ToLower(listener.Program)
			if strings.Contains(lowerListenerName, strings.ToLower(query)) {
				slog.Debug("Found listener", "program", listener.Program, "host", report.Host.HostName)
				rpt.Listeners = append(rpt.Listeners, listener)
				found = true
			}
		}
		if found {
			*reply = append(*reply, rpt)
		}
	}

	return nil
}
