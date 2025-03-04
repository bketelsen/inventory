package service

import (
	"errors"
	"fmt"
	"log"

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
