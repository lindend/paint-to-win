package communication

import (
	"paintToWin/gameserver/network"
)

type Message struct {
	Data       []byte
	Connection network.Connection
}

type InMessage struct {
	Data       []byte
	Connection network.Connection
	Entity     interface{}
}
