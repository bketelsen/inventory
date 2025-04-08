package inventory

// Storage defines the interface for persisting inventory data
type Storage interface {
	StoreReport(report Report) error
	GetReport(hostname string) (Report, bool)
	GetAllReports() []Report
}
