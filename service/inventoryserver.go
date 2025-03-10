package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/inventory/types"
)

// Define a struct. This struct will bind all the RPC methods
type InventoryServer struct {
	storage storage.Storage
	config  types.Config
}

// NewInventoryServer creates a new server with the provided storage
func NewInventoryServer(config types.Config, storage storage.Storage) *InventoryServer {
	return &InventoryServer{
		storage: storage,
		config:  config,
	}
}

// Update is the method that the client will call
func (c *InventoryServer) Update(report *types.Report, reply *int) error {
	if report == nil {
		return errors.New("report cannot be nil")
	}

	slog.Info("Received report", "host", report.Host.HostName, "ip", report.Host.IP)
	for _, l := range report.Listeners {
		slog.Info("Listener", "address", l.ListenAddress, "pid", l.PID, "port", l.Port, "program", l.Program, "app", l.Program)
	}
	for _, ct := range report.Containers {
		slog.Info("Container", "container", ct.ContainerID, "ip", ct.IP.String(), "hostname", ct.HostName, "platform", ct.Platform, "image", ct.Image)
	}
	for _, s := range report.Services {
		slog.Info("Service", "name", s.Name, "port", s.Port, "protocol", s.Protocol)
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

// Update is the method that the client will call
func (c *InventoryServer) Search(query string, reply *[]types.Report) error {
	if query == "" {
		return errors.New("query cannot be nil")
	}

	slog.Info("Received search query", "query", query)
	// Perform the search in persistent storage
	reports := c.storage.GetAllReports()
	if reports == nil {
		return errors.New("no reports found")
	}
	for _, report := range reports {
		found := false
		rpt := types.Report{
			Host:       report.Host,
			Services:   []types.Service{},
			Containers: []types.Container{},
			Listeners:  []types.Listener{},
		}
		// iterate through reports and find matches on services and containers
		for _, service := range report.Services {
			lowerServiceName := strings.ToLower(service.Name)
			if strings.Contains(lowerServiceName, strings.ToLower(query)) {
				slog.Info("Found service", "service", service.Name, "host", report.Host.HostName)
				rpt.Services = append(rpt.Services, service)
				found = true
			}
		}
		for _, container := range report.Containers {
			lowerContainerName := strings.ToLower(container.ContainerID)
			lowerImageName := strings.ToLower(container.Image)
			lowerHostName := strings.ToLower(container.HostName)
			if strings.Contains(lowerContainerName, strings.ToLower(query)) || strings.Contains(lowerImageName, strings.ToLower(query)) || strings.Contains(lowerHostName, strings.ToLower(query)) {
				slog.Info("Found container", "container", container.ContainerID, "host", report.Host.HostName)
				rpt.Containers = append(rpt.Containers, container)
				found = true
			}
		}
		for _, listener := range report.Listeners {

			lowerListenerName := strings.ToLower(listener.Program)
			if strings.Contains(lowerListenerName, strings.ToLower(query)) {
				slog.Info("Found listener", "program", listener.Program, "host", report.Host.HostName)
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
