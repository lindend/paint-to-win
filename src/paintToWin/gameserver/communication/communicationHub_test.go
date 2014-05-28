package communication

import (
	"paintToWin/gameserver/network"
	"testing"
	"time"
)

type playerMock struct {
}

type gameMock struct {
}

type connectionMock struct {
}

func (cmock connectionMock) Send(data []byte) error {
	return nil
}

func (cmock connectionMock) Close() error {
	return nil
}

func SetupTestHub() (*CommunicationHub, chan network.Connection, chan network.Connection, chan Message, <-chan Message) {
	connects := make(chan network.Connection)
	disconnects := make(chan network.Connection)
	inData := make(chan Message)

	hub, outData := NewCommunicationHub(connects, disconnects, inData, nil)
	go hub.Serve()
	return &hub, connects, disconnects, inData, outData
}

func TestRegisterChannel(t *testing.T) {
	hub, _, _, _, _ := SetupTestHub()
	_ = hub.RegisterChannel("testGame")

	if _, exists := hub.channels["testGame"]; !exists {
		t.Error("Channel was not registered on hub")
	}
}

func TestUnregisterChannel(t *testing.T) {
	hub, _, _, _, _ := SetupTestHub()
	_ = hub.RegisterChannel("testGame")
	hub.UnregisterChannel("testGame")
	if _, exists := hub.channels["testGame"]; exists {
		t.Error("Channel was not removed from hub")
	}
}

func TestNewConnectionCallsHandshake(t *testing.T) {
	hub, connects, _, _, _ := SetupTestHub()
	handshakeDone := make(chan bool)
	hub.handshake = func(conn network.Connection, inData <-chan InMessage, outData chan<- Message) (interface{}, string) {
		handshakeDone <- true
		return nil, ""
	}

	connects <- nil

	select {
	case <-time.After(50):
		t.Error("Handshake call timed out")
	case <-handshakeDone:
		//Was successful
	}
}

func TestNullHandshake(t *testing.T) {
	_, connects, _, _, _ := SetupTestHub()
	connects <- nil
}

func TestForwardsDataOnInChannel(t *testing.T) {
	hub, connects, _, inData, _ := SetupTestHub()
	gameChannel := hub.RegisterChannel("testGame")
	hub.handshake = func(conn network.Connection, inData <-chan InMessage, outData chan<- Message) (interface{}, string) {
		player := playerMock{}
		return player, "testGame"
	}
	conn := connectionMock{}
	connects <- conn
	inData <- Message{[]byte("TestData"), conn}
	result := <-gameChannel

	if string(result.Data) != "TestData" {
		t.Error("Invalid data forwarded on channel")
	}
}

func TestForwardsDataOnHandshakeChannel(t *testing.T) {
	hub, connects, _, inData, outData := SetupTestHub()
	hub.handshake = func(conn network.Connection, handshakeInData <-chan InMessage, handshakeOutData chan<- Message) (interface{}, string) {
		player := playerMock{}
		message := <-handshakeInData
		handshakeOutData <- Message{message.Data, message.Connection}

		return player, ""
	}
	conn := connectionMock{}
	connects <- conn
	inData <- Message{[]byte("TestMessage"), conn}
	message := <-outData

	if string(message.Data) != "TestMessage" {
		t.Error("Invalid handshake message received")
	}
}
