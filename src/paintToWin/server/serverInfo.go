package server

import (
	"net"
	"os"
)

type ServerInfo struct {
	Name     string
	HostName string
}

func LoadServerInfo() (ServerInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return ServerInfo{}, err
	}

	return ServerInfo{
		Name:     hostname,
		HostName: hostname,
	}, nil
}

func findIpAddress(hostname string) (string, error) {
	addresses, err := net.LookupHost(hostname)
	if err != nil {
		return "", err
	}
	return addresses[0], nil
}
