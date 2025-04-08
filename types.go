// Package inventory provides the data structures and types used in the inventory service
package inventory

import (
	"net"
	"time"
)

// ContainerPlatform represents the platform on which a container is running
type ContainerPlatform int

const (
	Docker ContainerPlatform = iota
	Incus
	KVM
	Podman
)

// String returns the string representation of the ContainerPlatform
func (p ContainerPlatform) String() string {
	switch p {
	case Docker:
		return "Docker"
	case Incus:
		return "Incus"
	case KVM:
		return "KVM"
	case Podman:
		return "Podman"
	default:
		return "Unknown"
	}
}

// Report represents the data structure sent to the server
// It contains information about the host, services, listeners, and containers
// that are running on the host.
type Report struct {
	Host       Host        `table:"host,default_sort"`
	Services   []Service   `table:"services"`
	Listeners  []Listener  `table:"listeners"`
	Containers []Container `table:"containers"`
	Timestamp  time.Time   `table:"timestamp" `
}

// DisplayTime returns the timestamp in a human-readable format
func (r *Report) DisplayTime() string {
	return r.Timestamp.Format(time.UnixDate)
}

// Host represents the host information
// It contains the hostname, IP address, location, and description of the host.
// The IP address is represented as a string to allow for both IPv4 and IPv6 addresses.
// The location and description fields are optional and can be used to provide
// additional information about the host.
// The location field can be used to specify the physical location of the host,
// such as a data center or office location.
// The description field can be used to provide a brief description of the host,
// such as its purpose or role in the network or a physical description of the host.
type Host struct {
	HostName    string `table:"hostname,default_sort" json:"host_name,omitempty"`
	IP          string `table:"ip" json:"ip,omitempty"`
	Location    string `table:"location" json:"location,omitempty"`
	Description string `table:"description" json:"description,omitempty"`
}

// Service represents a service running on the host
// It contains the name of the service, the port it is listening on,
// the address it is listening on, the protocol it uses, and the unit name
// (if applicable).
// Services can be any type of service, such as a web server, database server,
// or any other type of service that listens for incoming connections.
// They are provided to account for services that are not necessarily containers.
type Service struct {
	Name      string    `table:"name,default_sort" json:"name,omitempty"`
	Port      uint16    `table:"port" json:"port,omitempty"`
	Listeners []*Listen `table:"listeners" json:"listeners,omitempty"`
	Protocol  string    `table:"protocol" json:"protocol,omitempty"`
	Unit      string    `table:"unit" json:"unit,omitempty"`
}

// Listen represents a network listener for a service
type Listen struct {
	Port          uint16 `table:"port" json:"port,omitempty"`
	ListenAddress string `yaml:"listen_address" table:"listen_address" json:"listen_address,omitempty"`
	Protocol      string `table:"protocol" json:"protocol,omitempty"`
}

// Listener represents a network listener on the host
// It contains the address it is listening on, the port it is listening on,
// the PID of the process that is listening, and the name of the program
// that is listening.
type Listener struct {
	ListenAddress net.IP `table:"listen_address" json:"listen_address,omitempty"`
	Port          uint16 `table:"port" json:"port,omitempty"`
	PID           int    `table:"pid" json:"pid,omitempty"`
	Program       string `table:"program,default_sort" json:"program,omitempty"`
}

// Container represents a container running on the host
// It contains the container ID, image name, IP address, ports it is listening on,
// the hostname of the container, and the platform it is running on.
// The container ID is a unique identifier for the container,
// which can be used to manage or interact with the container.
// The image name is the name of the container image that was used to create the container.
// The IP address is the address assigned to the container,
// which can be used to communicate with the container.
// The ports field is a list of ports that the container is listening on,
// which can be used to access services running inside the container.
// The hostname is the name of the container, which can be used to identify the container
// in a network or for logging purposes.
// The platform field indicates the container platform that is being used,
// such as Docker, Incus, KVM, or Podman.
type Container struct {
	ContainerID string            `table:"container_id,default_sort" json:"container_id,omitempty"`
	Image       string            `table:"image" json:"image,omitempty"`
	IP          net.IPAddr        `table:"ip" json:"ip,omitempty"`
	Ports       []string          `table:"ports" json:"ports,omitempty"`
	HostName    string            `table:"hostname" json:"host_name,omitempty"`
	Platform    ContainerPlatform `table:"platform" json:"platform,omitempty"`
}
