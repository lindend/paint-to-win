package gamemanager

import (
	"paintToWin/gameserver/game"
	"paintToWin/storage"
)

func ToStorageGame(g *game.Game, isActive bool, numPlayers int, server storage.Server) *storage.Game {
	return &storage.Game{
		GameId:     g.Id,
		Name:       g.Name,
		CreatedBy:  "",
		IsActive:   isActive,
		NumPlayers: numPlayers,
		HostedOn:   server.Address,
	}
}
