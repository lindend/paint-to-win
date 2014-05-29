package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/web"
)

type CreateGameInput struct {
}

type CreateGameOutput struct {
	GameId    string
	Endpoints []network.EndpointInfo
}

type ReservationInput struct {
}

type ReservationOutput struct {
	GameId        string
	ReservationId string
	TimeToLive    int
	Endpoints     []network.EndpointInfo
}

func NewReservationOutput(gameId string, reservationId string, endpoints []network.EndpointInfo) ReservationOutput {
	return ReservationOutput{
		GameId:        gameId,
		ReservationId: reservationId,
		Endpoints:     endpoints,
	}
}

func CreateHandler(gameManager *gamemanager.GameManager) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		g, err := gameManager.CreateGame()
		fmt.Println("Created game ", g, err)
		if err != nil {
			return nil, web.NewApiError(http.StatusInternalServerError, "")
		} else {
			return CreateGameOutput{
				GameId:    g.GameId,
				Endpoints: gameManager.Endpoints(),
			}, nil
		}
	}
}

func ReserveHandler(gameManager *gamemanager.GameManager) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		vars := mux.Vars(req)
		gameId := vars["gameId"]
		playerId := vars["playerId"]

		reservationId, err := gameManager.ReserveSpot(playerId, gameId)
		if err != nil {
			return nil, web.NewApiError(http.StatusInternalServerError, err.Error())
		} else {
			return NewReservationOutput(gameId, reservationId, gameManager.Endpoints()), nil
		}
	}
}

func RegisterGameManagerApi(router *mux.Router, gameManager *gamemanager.GameManager) {
	router.HandleFunc("/games/create", web.DefaultHandler(CreateHandler(gameManager))).Methods("POST", "OPTIONS")
	router.HandleFunc("/games/{gameId}/reserve/{playerId}", web.DefaultHandler(ReserveHandler(gameManager))).Methods("POST", "OPTIONS")
}
