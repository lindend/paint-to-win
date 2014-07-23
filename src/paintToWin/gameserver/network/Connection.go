package network

import ()

type Connection struct {
	InData  <-chan Packet
	OutData chan<- Packet
	Closed  chan struct{}
}

type NewConnection struct {
	Connection Connection
	Variables  map[string]string
}
