package main

import (
	"fmt"

	"paintToWin/gameserver/codec"
	"paintToWin/gameserver/game"
	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/storage"
)

func ClientHandshake(gameManager *gamemanager.GameManager, store *storage.Storage, idGenerator <-chan string, connection network.NewConnection) {
	sessionId := connection.Variables["sessionId"]
	reservationId := connection.Variables["reservationId"]

	fmt.Println("New player")
	fmt.Println("SessionId ", sessionId)
	fmt.Println("ReservationId ", reservationId)

	session := &storage.Session{}
	if err := store.GetFromCache("session:"+sessionId, session); err != nil {
		fmt.Println("Handshake: no such session (", err, ")")
		return
	}

	var playerName string
	if !session.Player.IsGuest {
		storagePlayer := storage.Player{}
		if err := store.FirstWhere(storage.Player{Id: session.Player.UserId}, &storagePlayer); err != nil {
			fmt.Println("Handshake: no such player")
			return
		}
		playerName = storagePlayer.UserName
	} else {
		playerName = session.Player.PlayerName
	}

	playerOutData := codec.NewGameOutDataEncoder(connection.Connection)
	player := game.NewPlayer(playerName, session.Player.IsGuest, reservationId, <-idGenerator, playerOutData)

	g, err := gameManager.ClaimSpot(reservationId, player)
	if err != nil {
		fmt.Println("Handshake: no such reservation")
		close(playerOutData)
		return
	}
	codec.NewGameInDataDecoder(connection.Connection.InData, player, g)
}
