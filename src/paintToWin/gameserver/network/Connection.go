package network

import ()

type Connection interface {
	Close() error
	Send(data []byte) error
}

type NewConnection struct {
	Connection Connection
	Variables  map[string]string
}
