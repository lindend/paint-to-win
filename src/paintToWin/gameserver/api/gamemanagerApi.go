package api

import (
	"fmt"

	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/service"
)

type CreateGameInput struct {
	Name     string
	Password string
	Wordlist string
}

type CreateGameOutput struct {
	GameId    string
	Endpoints []network.EndpointInfo
}

type ReservationInput struct {
	GameId string
}

type ReservationOutput struct {
	GameId        string
	ReservationId string
	TimeToLive    int
	Endpoints     []network.EndpointInfo
}

var CreateGameOperation = service.NewOperation(serviceName, "createGame", "games", "POST", CreateGameInput{}, CreateGameOutput{})
var ReserveSpotOperation = service.NewOperation(serviceName, "reserveSpot", "games/{GameId}/reserve", "POST", ReservationInput{}, ReservationOutput{})

func NewReservationOutput(gameId string, reservationId string, endpoints []network.EndpointInfo) ReservationOutput {
	return ReservationOutput{
		GameId:        gameId,
		ReservationId: reservationId,
		Endpoints:     endpoints,
	}
}

func registerCreateGameOperation(host service.Host, gameManager *gamemanager.GameManager) {
	host.Register(func(input CreateGameInput) (CreateGameOutput, error) {
		fmt.Println("Creating game with input", input)
		g, err := gameManager.CreateGame(input.Name, input.Wordlist)
		fmt.Println("Created game ", g, err)
		if err != nil {
			return CreateGameOutput{}, err
		} else {
			return CreateGameOutput{
				GameId:    g.GameId,
				Endpoints: gameManager.Endpoints(),
			}, nil
		}
	}, CreateGameOperation)
}

func registerReserveSpotOperation(host service.Host, gameManager *gamemanager.GameManager) {
	host.Register(func(input ReservationInput) (ReservationOutput, error) {
		reservationId, err := gameManager.ReserveSpot(input.GameId)
		if err != nil {
			return ReservationOutput{}, err
		} else {
			return NewReservationOutput(input.GameId, reservationId, gameManager.Endpoints()), nil
		}
	}, ReserveSpotOperation)
}

func RegisterGameManagerApi(serviceHost service.Host, gameManager *gamemanager.GameManager) {
	registerCreateGameOperation(serviceHost, gameManager)
	registerReserveSpotOperation(serviceHost, gameManager)
}
