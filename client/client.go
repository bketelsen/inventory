// Package client provides a client for the inventory server.
// It handles the connection to the server and sending inventory reports
package client

import (
	"log/slog"
	"net/rpc"
	"os"

	"github.com/bketelsen/inventory"

	"github.com/jackpal/gateway"
)

// Client is the client for the inventory server
type Client struct {
	services    []*inventory.Service
	remote      string
	location    string
	description string
}

// NewClient creates a new client for the inventory server
// `remote` is the address of the inventory server
// `location` is the location of the host
// `description` is the description of the host
// `services` is a list of services to report
func NewClient(remote, location, description string, services []*inventory.Service) *Client {
	return &Client{
		remote:      remote,
		location:    location,
		description: description,
		services:    services,
	}
}

// Send sends the inventory report to the server
// It connects to the server using the remote address, and
// collects the host information, listeners, and containers
func (r *Client) Send() error {
	slog.Info("starting inventory client", "remote", r.remote)
	cl, err := rpc.Dial("tcp", r.remote) // Connect to server using config
	if err != nil {
		slog.Error("Error connecting to server", "address", r.remote, "error", err)
		return err
	}
	defer cl.Close()

	// Prepare arguments
	report := inventory.Report{}
	host, err := GetHost()
	if err != nil {
		slog.Error("Error getting host", "error", err)
		return err
	}
	report.Host = host
	report.Host.Location = r.location
	report.Host.Description = r.description

	report.Listeners, err = GetListeners()
	if err != nil {
		slog.Error("Error getting listeners", "error", err)
		return err
	}
	dockerContainers, err := GetDockerContainers(host.IP)
	if err != nil {
		slog.Error("Error getting docker containers:", "error", err)
		return err
	}
	incusContainers, err := GetIncusContainers()
	if err != nil {
		slog.Error("Error getting incus containers:", "error", err)
		return err
	}
	report.Containers = append(report.Containers, dockerContainers...)
	report.Containers = append(report.Containers, incusContainers...)

	if r.services != nil {
		// Add services to report
		for _, s := range r.services {
			report.Services = append(report.Services, *s)
		}
	}

	var result int

	// Call Update method on the server
	err = cl.Call("InventoryServer.Update", report, &result)
	if err != nil {
		slog.Error("Error calling InventoryServer.Update", "error", err)
		return err
	}
	slog.Info("Inventory report sent successfully", "result", result)
	return nil
}

// Search sends a search query to the server and returns the results
// It connects to the server using the remote address and
// sends the query to the server.
func (r *Client) Search(query string) ([]inventory.Report, error) {
	cl, err := rpc.Dial("tcp", r.remote) // Connect to server using config
	if err != nil {
		slog.Error("Error connecting to server", "address", r.remote, "error", err)
		return []inventory.Report{}, err
	}
	defer cl.Close()

	var result []inventory.Report
	// Call Search method on the server
	err = cl.Call("InventoryServer.Search", query, &result)
	if err != nil {
		slog.Error("Error calling InventoryServer.Search", "error", err)
		return []inventory.Report{}, err
	}

	return result, nil
}

func GetHost() (inventory.Host, error) {
	slog.Info("Getting host information")
	h := inventory.Host{}
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error getting hostname", "error", err)
		return h, err
	}
	h.HostName = hostname

	ip, err := gateway.DiscoverInterface()
	if err != nil {
		slog.Error("Error discovering IP address", "error", err)
		return h, err
	}
	h.IP = ip.String()

	slog.Info("Host information", "hostname", h.HostName, "IP", h.IP)
	return h, nil
}
