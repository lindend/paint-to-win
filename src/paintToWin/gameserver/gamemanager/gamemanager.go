package gamemanager

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"paintToWin/gameserver/game"
	"paintToWin/gameserver/gamestate"
	"paintToWin/gameserver/network"
	"paintToWin/storage"
	"paintToWin/wordlist"
)

const ReservationTimeout = 60

var PlayerNotFoundError = errors.New("player not found")
var GameDoesNotExistError = errors.New("game does not exist")

type gameItem struct {
	game        *game.Game
	storageGame *storage.Game
}

type GameManager struct {
	syncLock *sync.Mutex

	idGenerator  <-chan string
	storage      *storage.Storage
	games        map[string]*gameItem
	reservations map[string]*gameItem
	server       storage.Server

	endpoints []network.EndpointInfo
}

func NewGameManager(
	idGenerator <-chan string,
	endpoints []network.EndpointInfo,
	store *storage.Storage,
	server storage.Server,
) *GameManager {
	gameManager := GameManager{
		syncLock:     &sync.Mutex{},
		games:        make(map[string]*gameItem),
		reservations: make(map[string]*gameItem),
		idGenerator:  idGenerator,
		storage:      store,
		endpoints:    endpoints,
		server:       server,
	}
	return &gameManager
}

func (gameManager *GameManager) CreateGame(name string, wordlistId int64) (*storage.Game, error) {
	fmt.Println("Creating new game")
	words := loadWords(gameManager.storage, wordlistId, 20)

	newGame := game.NewGame(<-gameManager.idGenerator, gamestate.NewInitRoundState(words))
	newGame.Name = name
	storageGame := ToStorageGame(newGame, true, 0, gameManager.server)
	gameManager.storage.Save(storageGame)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Fatal error while serving game", r, string(debug.Stack()))
			}

			storageGame.IsActive = false
			gameManager.storage.Save(storageGame)
			gameManager.removeGame(newGame)
		}()
		newGame.Run()
		gameManager.removeGame(newGame)
	}()

	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	gameManager.games[newGame.Id] = &gameItem{newGame, storageGame}
	return storageGame, nil
}

func (gameManager *GameManager) ReserveSpot(gameId string) (string, error) {
	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	g, exists := gameManager.games[gameId]
	if !exists {
		return "", GameDoesNotExistError
	}

	reservationId := <-gameManager.idGenerator + <-gameManager.idGenerator
	gameManager.reservations[reservationId] = g

	go func() {
		<-time.After(ReservationTimeout * time.Second)
		gameManager.syncLock.Lock()
		delete(gameManager.reservations, reservationId)
		gameManager.syncLock.Unlock()
	}()

	return reservationId, nil
}

func (gameManager *GameManager) ClaimSpot(reservationId string, player *game.Player) (*game.Game, error) {
	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	g, exists := gameManager.reservations[reservationId]
	if !exists {
		return nil, GameDoesNotExistError
	}
	delete(gameManager.reservations, reservationId)

	if _, ok := gameManager.games[g.game.Id]; !ok {
		return nil, GameDoesNotExistError
	}

	if !g.game.PlayerJoin(player) {
		return nil, GameDoesNotExistError
	}

	return g.game, nil
}

func (gameManager *GameManager) ReclaimSpot(gameId string, playerId string) (*game.Game, error) {
	g, exists := gameManager.games[gameId]
	if !exists {
		return nil, GameDoesNotExistError
	}

	_, exists = getPlayer(g.game, playerId)
	if !exists {
		return nil, PlayerNotFoundError
	}

	return g.game, nil
}

func (gameManager *GameManager) removeGame(g *game.Game) {
	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	delete(gameManager.games, g.Id)
}

func loadWords(store *storage.Storage, wordlistId int64, numWords int) []string {
	dbWordlist := storage.Wordlist{}
	store.Where(&storage.Wordlist{Id: wordlistId}, &dbWordlist)
	wordlistWords := storage.WordlistWords{}
	store.Where(&storage.WordlistWords{WordlistId: wordlistId}, &wordlistWords)
	wlist := wordlist.LoadWordlistFromStorage(dbWordlist, wordlistWords)
	return wlist.GetWords(numWords)
}

func (gameManager GameManager) Endpoints() []network.EndpointInfo {
	return gameManager.endpoints
}

func getPlayer(g *game.Game, playerId string) (*game.Player, bool) {
	for _, p := range g.Players {
		if p.PlayerId == playerId {
			return p, true
		}
	}
	return nil, false
}
