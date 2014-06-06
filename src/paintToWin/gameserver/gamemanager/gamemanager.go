package gamemanager

import (
	"errors"
	"sync"
	"time"

	"paintToWin/gameserver/codec"
	"paintToWin/gameserver/communication"
	"paintToWin/gameserver/game"
	"paintToWin/gameserver/gamestate"
	"paintToWin/gameserver/network"
	"paintToWin/storage"
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

	communicationHub *communication.CommunicationHub
	endpoints        []network.EndpointInfo
}

func NewGameManager(
	idGenerator <-chan string,
	endpoints []network.EndpointInfo,
	commHub *communication.CommunicationHub,
	store *storage.Storage,
	server storage.Server,
) *GameManager {
	gameManager := GameManager{
		syncLock:         &sync.Mutex{},
		games:            make(map[string]*gameItem),
		reservations:     make(map[string]*gameItem),
		idGenerator:      idGenerator,
		storage:          store,
		endpoints:        endpoints,
		communicationHub: commHub,
		server:           server,
	}
	return &gameManager
}

func (gameManager *GameManager) CreateGame() (*storage.Game, error) {
	newGame := game.NewGame(<-gameManager.idGenerator, gamestate.NewInitRoundState())
	storageGame := ToStorageGame(newGame, true, 0, gameManager.server)
	gameManager.storage.Save(storageGame)

	gameInput := gameManager.communicationHub.RegisterChannel(newGame.Id)
	codec.NewGameMessageDecoder(gameInput, newGame.InData)

	go func() {
		newGame.Run()
		storageGame.IsActive = false
		gameManager.storage.Save(storageGame)
		gameManager.removeGame(newGame)
	}()

	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	gameManager.games[newGame.Id] = &gameItem{newGame, storageGame}
	return storageGame, nil
}

func (gameManager *GameManager) ReserveSpot(playerName string, gameId string) (string, error) {
	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	g, exists := gameManager.games[gameId]
	if !exists {
		return "", GameDoesNotExistError
	}

	reservationId := <-gameManager.idGenerator + <-gameManager.idGenerator + playerName
	gameManager.reservations[reservationId] = g

	go func() {
		<-time.After(ReservationTimeout * time.Second)
		gameManager.syncLock.Lock()
		delete(gameManager.reservations, reservationId)
		gameManager.syncLock.Unlock()
	}()

	return reservationId, nil
}

func (gameManager *GameManager) ClaimSpot(reservationId string) (*game.Game, error) {
	gameManager.syncLock.Lock()
	defer gameManager.syncLock.Unlock()

	g, exists := gameManager.reservations[reservationId]
	if !exists {
		return nil, GameDoesNotExistError
	}
	delete(gameManager.reservations, reservationId)
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
