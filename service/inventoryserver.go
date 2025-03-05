package service

import (
	"errors"
	"fmt"
	"log"
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

	log.Printf("Received report for %s at %s", report.Host.HostName, report.Host.IP)
	for _, l := range report.Listeners {
		log.Printf("[%s] Listener: Address=%s PID=%d, Port=%d, App=%s", report.Host.HostName, l.ListenAddress, l.PID, l.Port, l.Program)
	}
	for _, c := range report.Containers {
		log.Printf("[%s] Container: Name=%s IP=%s HostName=%s Platform=%s Image=%s", report.Host.HostName, c.ContainerID, c.IP.String(), c.HostName, c.Platform, c.Image)
	}
	for _, c := range report.Services {
		log.Printf("[%s] Service: Name=%s ListenAddress=%s Port=%d Protocol=%s Unit=%s", report.Host.HostName, c.Name, c.ListenAddress, c.Port, c.Protocol, c.Unit)
	}

	// Store the report in persistent storage
	err := c.storage.StoreReport(*report)
	if err != nil {
		log.Printf("Failed to store report: %v", err)
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

	log.Printf("Received search query: %s", query)
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
				log.Printf("Found service %s on %s", service.Name, report.Host.HostName)
				rpt.Services = append(rpt.Services, service)
				found = true
			}
		}
		for _, container := range report.Containers {
			lowerContainerName := strings.ToLower(container.ContainerID)
			lowerImageName := strings.ToLower(container.Image)
			lowerHostName := strings.ToLower(container.HostName)
			if strings.Contains(lowerContainerName, strings.ToLower(query)) || strings.Contains(lowerImageName, strings.ToLower(query)) || strings.Contains(lowerHostName, strings.ToLower(query)) {
				log.Printf("Found container %s on %s", container.ContainerID, report.Host.HostName)
				rpt.Containers = append(rpt.Containers, container)
				found = true
			}
		}
		for _, listener := range report.Listeners {

			lowerListenerName := strings.ToLower(listener.Program)
			if strings.Contains(lowerListenerName, strings.ToLower(query)) {
				log.Printf("Found listener %s on %s", listener.Program, report.Host.HostName)
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
