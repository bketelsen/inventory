//go:build linux
// +build linux

package client

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"

	"github.com/bketelsen/inclient"
	"github.com/bketelsen/inventory"
	"github.com/cakturk/go-netstat/netstat"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	config "github.com/lxc/incus/v6/shared/cliconfig"
)

// GetListeners returns a list of network listeners on the host
// It uses the go-netstat library to get the list of TCP and UDP sockets
// and filters them to only include those in the LISTEN state.
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

// GetDockerContainers returns a list of docker containers running on the host
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
