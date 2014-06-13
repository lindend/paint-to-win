package communication

import (
	"fmt"
	"paintToWin/gameserver/network"
	"sync"
)

type HandshakeFunction func(connection network.NewConnection, inData <-chan InMessage, outData chan<- Message) (commEntity interface{}, channelId string)

type connectionData struct {
	inData     chan<- InMessage
	entity     interface{}
	disconnect chan<- interface{}
}

type CommunicationChannel struct {
	InData     chan InMessage
	Disconnect chan interface{}
}

type CommunicationHub struct {
	data    <-chan Message
	outData chan<- Message

	connects    <-chan network.NewConnection
	disconnects <-chan network.Connection
	handshake   HandshakeFunction

	channels     map[string]CommunicationChannel
	connections  map[network.Connection]connectionData
	channelsLock sync.Mutex
}

func NewCommunicationHub(connects <-chan network.NewConnection, disconnects <-chan network.Connection, data <-chan Message) (*CommunicationHub, <-chan Message) {
	outData := make(chan Message)
	return &CommunicationHub{
		data,
		outData,

		connects,
		disconnects,
		nil,

		make(map[string]CommunicationChannel),
		make(map[network.Connection]connectionData),
		sync.Mutex{},
	}, outData
}

func (communicationHub *CommunicationHub) Serve(handshake HandshakeFunction) {
	communicationHub.handshake = handshake
	fmt.Println("Serving in communication hub")
	for {
		select {
		case connection := <-communicationHub.connects:
			fmt.Println("new connection connectionhub.go")
			handshakeChannel := make(chan InMessage)
			communicationHub.connections[connection.Connection] = connectionData{handshakeChannel, nil, nil}
			go communicationHub.newConnection(connection, handshakeChannel)
		case connection := <-communicationHub.disconnects:
			fmt.Println("disconnect communicationhub.go", connection)
			if connectionData, ok := communicationHub.connections[connection]; ok {
				fmt.Println("Sending disconnect to channel", connectionData.disconnect)
				connectionData.disconnect <- connectionData.entity
				delete(communicationHub.connections, connection)
			}
		case inMessage := <-communicationHub.data:
			fmt.Println("new message communicationhub.go")
			if connData, ok := communicationHub.connections[inMessage.Connection]; ok {
				fmt.Println("Sending data (communicationHub.go) ", connData.inData)
				connData.inData <- InMessage{inMessage.Data, inMessage.Connection, connData.entity}
			}
		}
	}
}

func (communicationHub *CommunicationHub) newConnection(connection network.NewConnection, handshakeChannel chan InMessage) {
	if communicationHub.handshake != nil {
		fmt.Println("Beginning handshake")
		entity, channelId := communicationHub.handshake(connection, handshakeChannel, communicationHub.outData)
		if entity != nil {
			if ch, exists := communicationHub.channels[channelId]; exists {
				fmt.Println("handshake complete")
				communicationHub.connections[connection.Connection] = connectionData{ch.InData, entity, ch.Disconnect}
			}
		} else {
			connection.Connection.Close()
		}
	}
}

func (communicationHub *CommunicationHub) RegisterChannel(id string) CommunicationChannel {
	inData := make(chan InMessage)
	disconnect := make(chan interface{})
	channel := CommunicationChannel{inData, disconnect}
	communicationHub.channels[id] = channel
	return channel
}

func (communicationHub *CommunicationHub) UnregisterChannel(id string) {
	if channel, ok := communicationHub.channels[id]; ok {
		close(channel.InData)
		close(channel.Disconnect)
		delete(communicationHub.channels, id)
	}
}
