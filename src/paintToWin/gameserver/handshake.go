package main

import (
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

		session := &storage.Session{}
		if err := store.GetFromCache("session:"+sessionId, session); err != nil {
			return nil, ""
		}

		storagePlayer := storage.Player{}
		if err := store.FirstWhere(storage.Player{Id: session.Player.Id}, &storagePlayer); err != nil {
			return nil, ""
		}
		player := game.NewPlayer(storagePlayer.UserName, false, reservationId, <-idGenerator, playerOutData)

		if g, err := gameManager.ClaimSpot(reservationId); err != nil {
			g.PlayerJoin <- player
			return player, g.Id
		}

		return nil, ""
	}
}
