package server

import (
	"os"
)

type ServerInfo struct {
	Name     string
	HostName string
	Address  string
}

func LoadServerInfo() (ServerInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return ServerInfo{}, err
	}

	//ip, err := findIpAddress()

	return ServerInfo{
		Name:     hostname,
		HostName: hostname,
		Address:  hostname,
	}, nil
}

func findIpAddress() (string, error) {
	return "", nil
}
