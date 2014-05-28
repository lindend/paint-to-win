package network

import ()

type Packet struct {
	Data       []byte
	Connection Connection
}
