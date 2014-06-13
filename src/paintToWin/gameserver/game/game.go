package game

import (
	"errors"
	"fmt"
	"time"
)

var PlayerNotFoundError = errors.New("Player not found")

type messageHandler func(from *Player, message Message)

type Game struct {
	Id      string
	Name    string
	Score   map[*Player]int
	Players []*Player

	InData      chan InMessage
	PlayerJoin  chan *Player
	PlayerLeave chan *Player

	NumRounds    int
	CurrentRound int

	messageHandler *MessageHandler

	state     []GameState
	isRunning bool

	timeout <-chan time.Time
}

func NewGame(id string, state GameState) *Game {
	g := &Game{
		Id:      id,
		Score:   make(map[*Player]int),
		Players: make([]*Player, 0),

		InData:      make(chan InMessage),
		PlayerJoin:  make(chan *Player),
		PlayerLeave: make(chan *Player),

		NumRounds:    10,
		CurrentRound: 0,

		messageHandler: NewMessageHandler(),

		state:     []GameState{state},
		isRunning: true,
		timeout:   time.After(10 * time.Minute),
	}
	g.setupMessageHandlers()
	return g
}

func (game *Game) setupMessageHandlers() {
	game.messageHandler.Add(game.playerChat)
}

func (game *Game) PushState(state GameState) {
	fmt.Println("Pushing state ", state, state.Name())
	game.ActiveState().Deactivate()
	game.state = append(game.state, state)
	game.activateState()
}

func (game *Game) PopState() GameState {
	game.ActiveState().Deactivate()
	state := game.state[len(game.state)-1]
	game.state = game.state[:len(game.state)-1]
	fmt.Println("Popping state to ", game.ActiveState(), game.ActiveState().Name())
	game.activateState()
	return state
}

func (game *Game) SwapState(state GameState) {
	fmt.Println("Swapping state to ", state, state.Name())
	game.ActiveState().Deactivate()
	game.state[len(game.state)-1] = state
	game.activateState()
}

func (game *Game) activateState() {
	activeState := game.ActiveState()
	players := []MessagePlayer{}
	for _, player := range game.Players {
		if !player.HasLeft {
			players = append(players, MessagePlayer{
				Id:   player.TempId,
				Name: player.Name,
			})
		}
	}
	game.Broadcast(NewGameStateMessage(activeState.Name(), activeState, players, game.timeLeft()))
	activeState.Activate(game)
}

func (game *Game) ActiveState() GameState {
	return game.state[len(game.state)-1]
}

func (game *Game) Stop() {
	game.isRunning = false
}

func (game *Game) SetTimeout(timeout time.Duration) {
	fmt.Println("New timeout ", timeout)
	game.timeout = time.After(timeout)
}

func (game *Game) timeLeft() int {
	return 0
}

func (game *Game) addPlayer(player *Player) {
	players := make([]MessagePlayer, 0, len(game.Players))
	msgPlayer := MessagePlayer{
		Id:   player.TempId,
		Name: player.Name,
	}

	joinMessage := NewPlayerJoinMessage(msgPlayer)
	for _, p := range game.Players {
		p.OutData <- joinMessage
		players = append(players, MessagePlayer{
			Id:   p.TempId,
			Name: p.Name,
		})
	}

	activeState := game.ActiveState()

	player.OutData <- NewWelcomeMessage(msgPlayer, "")
	player.OutData <- NewGameStateMessage(activeState.Name(), activeState, players, game.timeLeft())
	game.Players = append(game.Players, player)
}

func (game *Game) RemovePlayer(player *Player) {
	for i, pl := range game.Players {
		if pl == player {
			game.Players = append(game.Players[:i], game.Players[i+1:]...)
			break
		}
	}
	game.Broadcast(NewPlayerLeaveMessage(player.Name))
}

func (game *Game) NextPlayer(player *Player) *Player {
	return game.nextPlayerInGame(player, 1)
}

func (game *Game) PreviousPlayer(player *Player) *Player {
	return game.nextPlayerInGame(player, -1)
}

func (game *Game) nextPlayerInGame(player *Player, step int) *Player {
	if playerIndex, err := game.findPlayer(player); err == nil {
		i := playerIndex + step
		for i != playerIndex {
			if i < 0 {
				i = len(game.Players) - 1
			} else if i == len(game.Players) {
				i = 0
			}
			if !game.Players[i].HasLeft {
				return game.Players[i]
			}
		}
	}
	return nil
}

func (game *Game) findPlayer(player *Player) (int, error) {
	for i, pl := range game.Players {
		if pl == player {
			return i, nil
		}
	}
	return 0, PlayerNotFoundError
}

func (game *Game) Broadcast(message Message) {
	for _, p := range game.Players {
		p.OutData <- message
	}
}

func (game *Game) Strokes(from *Player, strokes []Stroke) {
	strokesMessage := NewStrokesMessage(strokes)
	for _, p := range game.Players {
		if p != from {
			p.OutData <- strokesMessage
		}
	}
}

func (game *Game) Run() {
	game.ActiveState().Activate(game)
	for len(game.state) > 0 && game.isRunning {
		select {
		case <-game.timeout:
			fmt.Println("Timeout occurred")
			game.ActiveState().Timeout()
		case player := <-game.PlayerJoin:
			game.addPlayer(player)
			game.ActiveState().PlayerJoin(player)
		case player := <-game.PlayerLeave:
			fmt.Println("Player left (game)")
			player.HasLeft = true
			game.ActiveState().PlayerLeave(player)
			game.RemovePlayer(player)
		case message, ok := <-game.InData:
			if !ok {
				break
			} else if !game.handleMessage(message) {
				game.ActiveState().Message(message)
			}
		}
	}
	fmt.Println("Game stopped")
}

func (game *Game) playerChat(player *Player, chatMessage *ChatMessage) {
	outChat := NewChatMessage(player.Name, "", chatMessage.Message)

	to := chatMessage.To
	for _, toPlayer := range game.Players {
		if toPlayer.Name == to || to == "" {
			toPlayer.OutData <- outChat
		}
	}
}

func (game *Game) handleMessage(message InMessage) bool {
	err := game.messageHandler.Handle(message)
	return err == nil
}
