package types

import (
	"net"
)

type ContainerPlatform int

const (
	Docker ContainerPlatform = iota
	Incus
	KVM
	Podman
)

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

type Report struct {
	Host       Host
	Services   []Service
	Listeners  []Listener
	Containers []Container
}

type Host struct {
	HostName    string
	IP          string
	Location    string
	Description string
}

type Service struct {
	Name          string
	Port          uint16
	ListenAddress string
	Protocol      string
	Unit          string
}

type Listener struct {
	ListenAddress net.IP
	Port          uint16
	PID           int
	Program       string
}

type Container struct {
	ContainerID string
	Image       string
	IP          net.IPAddr
	Ports       []string
	HostName    string
	Platform    ContainerPlatform
}
