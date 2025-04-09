// Package client provides a client for the inventory server.
// It handles the connection to the server and sending inventory reports
package client

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/bketelsen/inclient"
	"github.com/bketelsen/inventory"
	"github.com/cakturk/go-netstat/netstat"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jackpal/gateway"
	config "github.com/lxc/incus/v6/shared/cliconfig"
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

func GetListeners() ([]inventory.Listener, error) {
	slog.Info("Getting network listeners")

	l := []inventory.Listener{}
	// get tcp sockets
	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		slog.Error("Error getting TCP4 sockets", "error", err)
		return l, err
	}
	for _, tab := range tabs {
		ipaddr := tab.LocalAddr.IP
		process := "-"
		pid := 0
		if tab.Process != nil {
			process = tab.Process.Name
			pid = tab.Process.Pid
		}
		l = append(l, inventory.Listener{
			ListenAddress: ipaddr,
			PID:           pid,
			Program:       process,
			Port:          tab.LocalAddr.Port,
		})
	}
	// get tcp6 sockets
	tabs, err = netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		slog.Error("Error getting TCP6 sockets", "error", err)

		return l, err
	}
	for _, tab := range tabs {
		ipaddr := tab.LocalAddr.IP
		process := "-"
		pid := 0
		if tab.Process != nil {
			process = tab.Process.Name
			pid = tab.Process.Pid
		}
		l = append(l, inventory.Listener{
			ListenAddress: ipaddr,
			PID:           pid,
			Program:       process,
			Port:          tab.LocalAddr.Port,
		})
	}
	return l, nil
}

func GetDockerContainers(defaultIP string) ([]inventory.Container, error) {
	containers := []inventory.Container{}
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Error creating Docker client", "error", err)
		return containers, err
	}
	defer apiClient.Close()

	cc, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: false})
	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			slog.Error("Docker not running, or can't connect", "error", err)
			return containers, nil
		}
		slog.Error("Error listing Docker containers", "error", err)
		return containers, err
	}

	for _, ctr := range cc {
		ports := []string{}
		log.Printf("Found Docker container: %s, Image: %s", ctr.ID[0:12], ctr.Image)
		for _, port := range ctr.Ports {
			if port.PublicPort == 0 {
				continue
			}
			if port.IP == "0.0.0.0" {
				port.IP = defaultIP
			}
			slog.Info("Found published port", "ip", port.IP, "publicPort", port.PublicPort, "type", port.Type)
			ports = append(ports, fmt.Sprintf("%s:%d/%s", port.IP, port.PublicPort, port.Type))
		}

		containers = append(containers, inventory.Container{
			ContainerID: ctr.ID[0:12],
			Image:       ctr.Image,
			Ports:       ports,
			HostName:    strings.TrimPrefix(ctr.Names[0], "/"),
			Platform:    inventory.Docker,
		})
	}
	return containers, nil
}

// GetIncusContainers returns a list of incus containers running on the host
func GetIncusContainers() ([]inventory.Container, error) {
	cfg := config.NewConfig("", true)
	containers := []inventory.Container{}
	client, err := inclient.NewClient(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "appear to be started") {
			slog.Error("No incus server running, or can't connect", "error", err)
			return containers, nil
		}
		slog.Error("Error creating Incus client", "error", err)
		return containers, err
	}
	cc, err := client.Instances(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "appear to be started") {
			slog.Error("No incus server running, or can't connect", "error", err)
			return containers, nil
		}
		slog.Error("Error listing Incus instances", "error", err)
		return nil, err
	}
	for _, c := range cc {
		if c.Status == "Stopped" {
			continue
		}
		if c.Status == "Error" {
			slog.Error("Incus container in error state", "name", c.Name)
			continue
		}
		slog.Info("Found Incus container", "name", c.Name, "status", c.Status)
		containers = append(containers, inventory.Container{
			ContainerID: c.Name,
			HostName:    c.Name,
			Image:       c.Config["image.os"] + "/" + c.Config["image.release"],
			Platform:    inventory.Incus,
			IP:          net.IPAddr{IP: net.ParseIP(c.State.Network["eth0"].Addresses[0].Address)},
		},
		)
	}
	return containers, nil
}
