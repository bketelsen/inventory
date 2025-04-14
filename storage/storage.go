// Package storage provides an in-memory implementation of the inventory.Storage interface
// It is the only storage implementation currently available
package storage

import (
	"log/slog"
	"slices"
	"strings"
	"sync"

	"github.com/bketelsen/inventory"
)

// ensure that MemoryStorage implements the Storage interface
var _ inventory.Storage = (*MemoryStorage)(nil)

// MemoryStorage implements the Storage interface with an in-memory map
type MemoryStorage struct {
	mu      sync.RWMutex
	reports map[string]inventory.Report
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	slog.Info("Creating new in-memory storage")
	return &MemoryStorage{
		reports: make(map[string]inventory.Report),
	}
}

// StoreReport stores a report in memory, keyed by hostname
func (ms *MemoryStorage) StoreReport(report inventory.Report) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if report.Host.HostName == "" {
		return inventory.ErrNoHostName
	}
	slog.Info("Storing report", "host", report.Host.HostName)
	ms.reports[report.Host.HostName] = report
	return nil
}

// GetReport retrieves a report by hostname
func (ms *MemoryStorage) GetReport(hostname string) (inventory.Report, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	slog.Info("Retrieving report", "host", hostname)
	report, exists := ms.reports[hostname]
	return report, exists
}

// GetAllReports returns all stored reports
func (ms *MemoryStorage) GetAllReports() []inventory.Report {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	slog.Info("Retrieving reports", "count", len(ms.reports))
	reports := make([]inventory.Report, 0, len(ms.reports))
	for _, report := range ms.reports {
		reports = append(reports, report)
	}
	// sort the reports by hostname
	slices.SortFunc(reports, func(a, b inventory.Report) int {
		if n := strings.Compare(a.Host.HostName, b.Host.HostName); n != 0 {
			return n
		}
		return 0
	})
	return reports
}
