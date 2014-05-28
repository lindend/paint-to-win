package game

import (
	"errors"
	"fmt"
	"math/rand"
	"paintToWin/gameserver/api"
	"paintToWin/storage"
	"paintToWin/web"
)

var NoActiveGameServersError = errors.New("No active game servers found")

func GetActiveGames(store *storage.Storage) ([]storage.Game, error) {
	var activeGames []storage.Game
	err := store.Where(&storage.Game{IsActive: true}, &activeGames)
	return activeGames, err
}

func removeGameServer(server *storage.Server) {

}

func CreateGame(store *storage.Storage) (storage.Game, error) {
	var gameServers []storage.Server
	if err := store.Where(&storage.Server{Type: "gameserver"}, &gameServers); err != nil {
		return storage.Game{}, err
	}
	if len(gameServers) == 0 {
		return storage.Game{}, NoActiveGameServersError
	}
	serverIndex := rand.Intn(len(gameServers))
	gameServer := gameServers[serverIndex]

	result := api.CreateGameOutput{}
	var errResult string
	fmt.Println("Creating game on ", gameServer.Address+"/games/create")
	if err := web.Post(gameServer.Address+"/games/create", nil, &result, &errResult); err != nil {
		fmt.Println("Error in http request ", err)
		return storage.Game{}, err
	}

	return storage.Game{}, nil
}

func JoinGame(gameId string, store *storage.Storage, session *storage.Session) (api.ReservationOutput, error) {
	game := &storage.Game{}
	if err := store.FirstWhere(storage.Game{GameId: gameId}, game); err != nil {
		return api.ReservationOutput{}, err
	}

	playerId := session.Player.Id

	result := api.ReservationOutput{}
	var errResult string
	if err := web.Post(fmt.Sprintf("%v/games/%v/reserve/%v", game.HostedOn,
		game.GameId, playerId), nil, &result, &errResult); err != nil {
		return api.ReservationOutput{}, err
	}

	return result, nil
}
