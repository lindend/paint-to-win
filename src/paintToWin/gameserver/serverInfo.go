package main

import (
	"os"
)

type serverInfo struct {
	Name     string
	HostName string
	Address  string
}

func loadServerInfo() (serverInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return serverInfo{}, err
	}

	return serverInfo{
		Name:     hostname,
		HostName: hostname,
		Address:  hostname,
	}, nil
}
