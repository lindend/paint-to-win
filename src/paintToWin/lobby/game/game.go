package game

import (
	"errors"
	"fmt"
	"math/rand"

	"paintToWin/gameserver/api"
	"paintToWin/service"
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

func CreateGame(services service.ServiceManager, name string, password string, wordlistId string) (storage.Game, error) {
	gameServers, err := services.Find("gameserver")
	if err != nil {
		return storage.Game{}, err
	}

	if len(gameServers) == 0 {
		return storage.Game{}, NoActiveGameServersError
	}

	for len(gameServers) > 0 {
		serverIndex := rand.Intn(len(gameServers))
		gameServer := gameServers[serverIndex]
		gameServers = append(gameServers[:serverIndex], gameServers[serverIndex+1:]...)

		apiInput := api.CreateGameInput{
			Name:     name,
			Password: password,
			Wordlist: wordlistId,
		}

		fmt.Println("Creating game on ", gameServer.Address+"/games")

		result, err := service.Call(api.CreateGameOperation, gameServer, apiInput)
		if err == nil {
			apiRes := result.(api.CreateGameOutput)
			return storage.Game{GameId: apiRes.GameId, Name: name}, nil
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
	if err := web.Post(fmt.Sprintf("%v/games/%v/reserve", game.HostedOn, game.GameId),
		nil, &result, &errResult); err != nil {
		return api.ReservationOutput{}, err
	}

	return result, nil
}

func GetWordLists(store *storage.Storage) ([]storage.Wordlist, error) {
	var wordLists []storage.Wordlist
	err := store.Db.Find(&wordLists).Error
	return wordLists, err
}
