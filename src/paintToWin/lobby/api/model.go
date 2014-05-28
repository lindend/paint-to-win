package api

import (
	"paintToWin/gameserver/network"
	"paintToWin/storage"
)

type Game struct {
	GameId     string
	Name       string
	CreatedBy  string
	IsActive   bool
	NumPlayers int
}

func NewGame(game storage.Game) Game {
	return Game{
		GameId:     game.GameId,
		Name:       game.Name,
		CreatedBy:  game.CreatedBy,
		IsActive:   game.IsActive,
		NumPlayers: game.NumPlayers,
	}
}

type Session struct {
	SessionId   string
	DisplayName string
}

func NewSession(sessionId string, displayName string) Session {
	return Session{
		SessionId:   sessionId,
		DisplayName: displayName,
	}
}

type SlotReservation struct {
	GameId        string
	ReservationId string
	Endpoints     []network.EndpointInfo
}
