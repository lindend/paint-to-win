package main

import (
	"fmt"

	"paintToWin/gameserver/codec"
	"paintToWin/gameserver/communication"
	"paintToWin/gameserver/game"
	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/storage"
)

func CreateClientHandshake(gameManager *gamemanager.GameManager, store *storage.Storage, idGenerator <-chan string) communication.HandshakeFunction {
	return func(connection network.NewConnection, inData <-chan communication.InMessage, outData chan<- communication.Message) (commEntity interface{}, channelId string) {
		playerOutData := make(chan game.Message)
		codec.NewGameMessageEncoder(playerOutData, outData, connection.Connection)

		sessionId := connection.Variables["sessionId"]
		reservationId := connection.Variables["reservationId"]

		fmt.Println("New player")
		fmt.Println("SessionId ", sessionId)
		fmt.Println("ReservationId ", reservationId)

		session := &storage.Session{}
		if err := store.GetFromCache("session:"+sessionId, session); err != nil {
			fmt.Println("Handshake: no such session (", err, ")")
			return nil, ""
		}

		var playerName string
		if !session.Player.IsGuest {
			storagePlayer := storage.Player{}
			if err := store.FirstWhere(storage.Player{Id: session.Player.UserId}, &storagePlayer); err != nil {
				fmt.Println("Handshake: no such player")
				return nil, ""
			}
			playerName = storagePlayer.UserName
		} else {
			playerName = session.Player.PlayerName
		}

		player := game.NewPlayer(playerName, session.Player.IsGuest, reservationId, <-idGenerator, playerOutData)

		if g, err := gameManager.ClaimSpot(reservationId); err != nil {
			fmt.Println("Handshake: no such reservation")
			return nil, ""
		} else {
			g.PlayerJoin <- player
			return player, g.Id
		}
	}
}
