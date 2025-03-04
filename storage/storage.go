package storage

import (
	"log"
	"sync"

	"github.com/bketelsen/inventory/types"
)

// Storage defines the interface for persisting inventory data
type Storage interface {
	StoreReport(report types.Report) error
	GetReport(hostname string) (types.Report, bool)
	GetAllReports() []types.Report
}

// MemoryStorage implements the Storage interface with an in-memory map
type MemoryStorage struct {
	mu      sync.RWMutex
	reports map[string]types.Report
	config  types.Config
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage(config types.Config) *MemoryStorage {
	log.Println("Creating new in-memory storage")
	return &MemoryStorage{
		config:  config,
		reports: make(map[string]types.Report),
	}
}

// StoreReport stores a report in memory, keyed by hostname
func (ms *MemoryStorage) StoreReport(report types.Report) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	log.Printf("[%s] Storing report", report.Host.HostName)
	ms.reports[report.Host.HostName] = report
	return nil
}

// GetReport retrieves a report by hostname
func (ms *MemoryStorage) GetReport(hostname string) (types.Report, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	log.Printf("Retrieving report for host %s", hostname)
	report, exists := ms.reports[hostname]
	return report, exists
}

// GetAllReports returns all stored reports
func (ms *MemoryStorage) GetAllReports() []types.Report {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	log.Printf("Retrieving %d reports", len(ms.reports))
	reports := make([]types.Report, 0, len(ms.reports))
	for _, report := range ms.reports {
		reports = append(reports, report)
	}
	return reports
}
