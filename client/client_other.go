//go:build !linux
// +build !linux

package client

import (
	"log/slog"

	"github.com/bketelsen/inventory"
)

// GetListeners returns an empty list of listeners
// because we are not running on Linux
// TODO: implement this for other OSes
func GetListeners() ([]inventory.Listener, error) {
	slog.Info("Skipping network listeners")

	l := []inventory.Listener{}
	return l, nil
}

// GetDockerContainers returns an empty list of containers
// because we are not running on Linux
// TODO: implement this for other OSes? Since docker would be
// running in a VM, will people be using it for hosting?
func GetDockerContainers(defaultIP string) ([]inventory.Container, error) {
	containers := []inventory.Container{}
	return containers, nil
}

// GetIncusContainers returns a list of incus containers running on the host
// GetIncusContainers returns an empty list of containers
// because we are not running on Linux
func GetIncusContainers() ([]inventory.Container, error) {
	containers := []inventory.Container{}

	return containers, nil
}
