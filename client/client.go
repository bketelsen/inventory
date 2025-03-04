package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/bketelsen/inclient"
	"github.com/bketelsen/inventory/types"
	"github.com/cakturk/go-netstat/netstat"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jackpal/gateway"
	config "github.com/lxc/incus/v6/shared/cliconfig"
)

type Reporter struct {
}

func NewReporter() *Reporter {
	return &Reporter{}
}
func (r *Reporter) Send() error {
	// Read config
	config, err := types.ReadConfig()
	if err != nil {
		log.Println("Error reading config:", err)
		return err
	}

	log.Println("Verbose mode:", config.Verbose)
	log.Println("Server address:", config.Server.Address)
	cl, err := rpc.Dial("tcp", config.Server.Address) // Connect to server using config
	if err != nil {
		log.Printf("Error connecting to %s: %v", config.Server.Address, err)
		return err
	}
	defer cl.Close()

	// Prepare arguments
	report := types.Report{}
	host, err := GetHost()
	if err != nil {
		log.Println("Error getting host:", err)
		return err
	}
	report.Host = host
	report.Listeners, err = GetListeners()
	if err != nil {
		log.Println("Error getting listeners:", err)
		return err
	}
	dockerContainers, err := GetDockerContainers(host.IP)
	if err != nil {
		log.Println("Error getting docker containers:", err)
		return err
	}
	incusContainers, err := GetIncusContainers()
	if err != nil {
		log.Println("Error getting incus containers:", err)
		return err
	}
	report.Containers = append(report.Containers, dockerContainers...)
	report.Containers = append(report.Containers, incusContainers...)

	report.Services = config.Services
	var result int

	// Call Update method on the server
	err = cl.Call("InventoryServer.Update", report, &result)
	if err != nil {
		log.Println("Error calling InventoryServer.Update:", err)
		return err
	}

	log.Println("server response:", result)
	return nil

}

func GetHost() (types.Host, error) {
	log.Println("Getting host information")
	h := types.Host{}
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
		return h, err
	}
	h.HostName = hostname

	ip, err := gateway.DiscoverInterface()
	if err != nil {
		log.Printf("Error getting IP address: %v", err)
		return h, err
	}
	h.IP = ip.String()

	log.Printf("Host: %s, IP: %s", h.HostName, h.IP)
	return h, nil
}

func GetListeners() ([]types.Listener, error) {
	log.Println("Getting network listeners")

	l := []types.Listener{}
	// get tcp sockets
	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		log.Printf("Error getting IPv4 TCP sockets: %v", err)
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
		l = append(l, types.Listener{
			ListenAddress: ipaddr,
			PID:           pid,
			Program:       process,
			Port:          uint16(tab.LocalAddr.Port),
		})
	}
	// get tcp6 sockets
	tabs, err = netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		log.Printf("Error getting IPv6 TCP sockets: %v", err)

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
		l = append(l, types.Listener{
			ListenAddress: ipaddr,
			PID:           pid,
			Program:       process,
			Port:          uint16(tab.LocalAddr.Port),
		})
	}
	return l, nil
}

func GetDockerContainers(defaultIP string) ([]types.Container, error) {
	containers := []types.Container{}
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return containers, err

	}
	defer apiClient.Close()

	cc, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: false})
	if err != nil {
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
			log.Printf("Docker not running, or can't connect")
			return containers, nil
		} else {
			log.Printf("Docker error: %v", err)

			return containers, err
		}
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
			log.Printf("Found published port: %s:%d/%s\n", port.IP, port.PublicPort, port.Type)
			ports = append(ports, fmt.Sprintf("%s:%d/%s", port.IP, port.PublicPort, port.Type))
		}

		containers = append(containers, types.Container{
			ContainerID: ctr.ID[0:12],
			Image:       ctr.Image,
			Ports:       ports,
			HostName:    strings.TrimPrefix(ctr.Names[0], "/"),
			Platform:    types.Docker,
		})
	}
	return containers, nil
}

// GetIncusContainers returns a list of incus containers running on the host
func GetIncusContainers() ([]types.Container, error) {
	cfg := config.NewConfig("", true)
	containers := []types.Container{}
	client, err := inclient.NewClient(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "appear to be started") {
			log.Println("No incus server running, or can't connect")
			return containers, nil
		} else {
			log.Printf("Error creating Incus client: %v", err)
			return containers, err
		}
	}
	cc, err := client.Instances(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "appear to be started") {
			log.Println("No incus server running, or can't connect")

			return containers, nil
		} else {
			log.Printf("Error listing Incus instances: %v", err)

			return nil, err
		}
	}
	for _, c := range cc {
		if c.Status == "Stopped" {
			continue
		}
		if c.Status == "Error" {
			log.Println("incus container in error state:", c.Name)
			continue
		}
		log.Println("Found Incus container:", c.Name)
		containers = append(containers, types.Container{
			ContainerID: c.Name,
			HostName:    c.Name,
			Image:       c.Config["image.os"] + "/" + c.Config["image.release"],
			Platform:    types.Incus,
			IP:          net.IPAddr{IP: net.ParseIP(c.State.Network["eth0"].Addresses[0].Address)}},
		)

	}
	return containers, nil
}
