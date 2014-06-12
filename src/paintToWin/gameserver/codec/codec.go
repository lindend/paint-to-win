package codec

import (
	"fmt"
	"paintToWin/gameserver/communication"
	"paintToWin/gameserver/game"
	"paintToWin/gameserver/network"
)

func StandardDecoder(inData <-chan network.Packet) <-chan communication.Message {
	resultChan := make(chan communication.Message)

	go func() {
		for {
			if msg, ok := <-inData; ok {
				fmt.Println("Standard decoder data")
				resultChan <- communication.Message{msg.Data, msg.Connection}
			} else {
				close(resultChan)
				break
			}
		}
	}()

	return resultChan
}

func StandardEncoder(inData <-chan communication.Message) <-chan network.Packet {
	resultChan := make(chan network.Packet)

	go func() {
		for {
			if msg, ok := <-inData; ok {
				fmt.Println("Standard encoder data")
				resultChan <- network.Packet{msg.Data, msg.Connection}
			} else {
				close(resultChan)
				break
			}
		}
	}()

	return resultChan
}

func NewGameMessageEncoder(inData <-chan game.Message, outData chan<- communication.Message, connection network.Connection) {
	go func() {
		for {
			if msg, ok := <-inData; ok {
				fmt.Println("Packet input codec.go")
				messageData, err := game.EncodeMessage(msg)
				if err != nil {
					fmt.Println("Error while encoding message", err)
				} else {
					outData <- communication.Message{messageData, connection}
				}
			} else {
				close(outData)
				break
			}
		}
	}()
}

func NewGameMessageDecoder(inData <-chan communication.InMessage, outData chan<- game.InMessage) {
	go func() {
		for {
			if msg, ok := <-inData; ok {
				player := msg.Entity.(*game.Player)
				gameMsg, err := game.DecodeMessage(msg.Data)
				if err == nil {
					outData <- game.InMessage{gameMsg, player}
				} else {
					fmt.Println("Error while decoding packet ", err)
				}
			} else {
				close(outData)
				break
			}
		}
	}()
}
