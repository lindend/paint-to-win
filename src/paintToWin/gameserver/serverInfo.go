package main

import (
	"fmt"
	"os"
	"paintToWin/storage"
)

type ServerInfo struct {
	Name     string
	HostName string
}

func loadServerInfo(port int) (storage.Server, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return storage.Server{}, "", err
	}

	return storage.Server{
		Name:    hostname,
		Address: fmt.Sprintf("http://%v:%d", hostname, port),
		Type:    "gameserver",
	}, hostname, nil
}
