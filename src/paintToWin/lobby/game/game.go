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

func CreateGame(store *storage.Storage, name string, password string, wordlistId int64) (storage.Game, error) {
	var gameServers []storage.Server
	if err := store.Where(&storage.Server{Type: "gameserver"}, &gameServers); err != nil {
		return storage.Game{}, err
	}
	if len(gameServers) == 0 {
		return storage.Game{}, NoActiveGameServersError
	}
	fmt.Println("Found ", len(gameServers), " game servers")
	var err error
	for len(gameServers) > 0 {
		serverIndex := rand.Intn(len(gameServers))
		gameServer := gameServers[serverIndex]
		gameServers = append(gameServers[:serverIndex], gameServers[serverIndex+1:]...)

		result := api.CreateGameOutput{}
		var errResult string
		fmt.Println("Creating game on ", gameServer.Address+"/games/create")
		apiInput := api.CreateGameInput{
			Name:     name,
			Password: password,
			Wordlist: wordlistId,
		}
		if err = web.Post(gameServer.Address+"/games/create", apiInput, &result, &errResult); err != nil {
			fmt.Println("Error in http request ", err)
		} else {
			return storage.Game{GameId: result.GameId, Name: name}, nil
		}
	}

	return storage.Game{}, err
}

func JoinGame(gameId string, store *storage.Storage, session *storage.Session) (api.ReservationOutput, error) {
	game := &storage.Game{}
	if err := store.FirstWhere(storage.Game{GameId: gameId}, game); err != nil {
		return api.ReservationOutput{}, err
	}

	result := api.ReservationOutput{}
	var errResult string
	if err := web.Post(fmt.Sprintf("%v/games/%v/reserve", game.HostedOn,
		game.GameId), nil, &result, &errResult); err != nil {
		return api.ReservationOutput{}, err
	}

	return result, nil
}

func GetWordLists(store *storage.Storage) ([]storage.WordList, error) {
	var wordLists []storage.WordList
	err := store.Db.Find(&wordLists).Error
	return wordLists, err
}
