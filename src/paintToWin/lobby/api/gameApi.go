package api

import (
	"net/http"
	"paintToWin/storage"
	"paintToWin/web"

	"paintToWin/lobby/game"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type JoinGameInput struct {
	GameId string
}

type CreateGameInput struct {
	Name      string `json: "name"`
	IsPrivate bool   `json: "isPrivate"`
	Password  string `json: "password"`
}

func (input *CreateGameInput) Validate() []web.InputError {
	return []web.InputError{}
}

func ListGamesHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		activeGames, _ := game.GetActiveGames(store)
		outGames := []Game{}
		for _, g := range activeGames {
			outGames = append(outGames, NewGame(g))
		}
		return outGames, nil
	}
}

func CreateGameHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		var input CreateGameInput
		inputErrs, err := web.DeserializeAndValidateInput(req, &input)
		if err != nil {
			return nil, web.NewApiError(http.StatusBadRequest, err)
		} else if inputErrs != nil {
			return nil, web.NewApiError(http.StatusBadRequest, inputErrs)
		}

		createdGame, _ := game.CreateGame(store, input.Name, input.Password)
		return NewGame(createdGame), nil
	}
}

func GetGameHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		return nil, web.NewApiError(http.StatusNotFound, nil)
	}
}

func JoinGameHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		vars := mux.Vars(req)
		gameId := vars["gameId"]
		session := context.Get(req, "session").(*storage.Session)
		if reserveOutput, err := game.JoinGame(gameId, store, session); err != nil {
			return nil, web.NewApiError(http.StatusInternalServerError, err.Error())
		} else {
			return SlotReservation{
				GameId:        reserveOutput.GameId,
				ReservationId: reserveOutput.ReservationId,
				Endpoints:     reserveOutput.Endpoints,
			}, nil
		}
	}
}

func RegisterGameApi(router *mux.Router, store *storage.Storage) {
	authenticator := NewSessionAuthenticator(store)

	router.HandleFunc("/games", web.DefaultAuthenticateHandler(ListGamesHandler(store), authenticator)).Methods("GET", "OPTIONS")
	router.HandleFunc("/games/create", web.DefaultAuthenticateHandler(CreateGameHandler(store), authenticator)).Methods("POST", "OPTIONS")

	router.HandleFunc("/games/{gameId}", web.DefaultAuthenticateHandler(GetGameHandler(store), authenticator)).Methods("GET", "OPTIONS")
	router.HandleFunc("/games/{gameId}/join", web.DefaultAuthenticateHandler(JoinGameHandler(store), authenticator)).Methods("POST", "OPTIONS")
}
