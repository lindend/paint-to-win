package codec

import (
	"fmt"
	"paintToWin/gameserver/game"
	"paintToWin/gameserver/network"
)

func NewGameInDataDecoder(inData <-chan network.Packet, player *game.Player, g *game.Game) {
	go func() {
		for msg := range inData {
			gameMsg, err := game.DecodeMessage(msg.Data)
			if err != nil {
				fmt.Println("Error while decoding packet", err)
				continue
			}
			if !g.OnData(player, gameMsg) {
				return
			}
		}
		fmt.Println("GameInDataDecoder: player left", player.Name)
		g.PlayerLeft(player)
	}()
}

func NewGameOutDataEncoder(connection network.Connection) chan<- game.Message {
	gameOutChan := make(chan game.Message)
	go func() {
		for gameMsg := range gameOutChan {
			data, err := game.EncodeMessage(gameMsg)
			if err != nil {
				fmt.Println("Error while encoding packet", err)
			}
			pkt := network.Packet{
				Data:       data,
				Connection: connection,
			}
			select {
			case connection.OutData <- pkt:
			case <-connection.Closed:
			}
		}
		close(connection.OutData)
	}()
	return gameOutChan
}
