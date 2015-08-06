package game

import (
	"errors"
	"fmt"
	"time"
	"math"
)

var PlayerNotFoundError = errors.New("Player not found")

type inMessage struct {
	Message Message
	Source  *Player
}

type Game struct {
	Id      string
	Name    string
	score   map[*Player]int
	Players PlayerList

	inData      chan inMessage
	playerJoin  chan *Player
	playerLeave chan *Player
	closed      chan struct{}

	NumRounds    int
	CurrentRound int

	messageHandler *MessageHandler

	state     []GameState
	isRunning bool

	timeout <-chan time.Time
	timeoutTime time.Time
}

func NewGame(id string, state GameState) *Game {
	g := &Game{
		Id:      id,
		score:   make(map[*Player]int),
		Players: make([]*Player, 0),

		inData:      make(chan inMessage),
		playerJoin:  make(chan *Player),
		playerLeave: make(chan *Player),
		closed:      make(chan struct{}),

		NumRounds:    10,
		CurrentRound: 0,

		messageHandler: NewMessageHandler(),

		state:     []GameState{state},
		isRunning: true,
		timeout:   time.After(10 * time.Minute),
		timeoutTime: time.Now().Add(10 * time.Minute),
	}
	g.setupMessageHandlers()
	return g
}

func (game *Game) OnData(from *Player, message Message) bool {
	select {
	case game.inData <- inMessage{message, from}:
		return true
	case <-game.closed:
		return false
	}
}

func (game *Game) PlayerLeft(player *Player) bool {
	select {
	case game.playerLeave <- player:
		return true
	case <-game.closed:
		return false
	}
}

func (game *Game) PlayerJoin(player *Player) bool {
	select {
	case game.playerJoin <- player:
		return true
	case <-game.closed:
		return false
	}
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
	activeState.Activate(game)
}

func (game *Game) BroadcastActiveState() {
	activeState := game.ActiveState()
	players := []MessagePlayer{}
	for _, player := range game.Players {
		if !player.HasLeft {
			players = append(players, MessagePlayer{
				Id:      player.TempId,
				Name:    player.Name,
				IsGuest: player.IsGuest,
				Score:   game.score[player],
			})
		}
	}
	game.Broadcast(NewGameStateMessage(activeState.Name(), activeState, players, game.timeLeft()))
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
	game.timeoutTime = time.Now().Add(timeout)
}

func (g *Game) Broadcast(message Message) {
	for _, p := range g.Players {
		p.OutData <- message
	}
}

func (game *Game) timeLeft() int {
	timeDifference := game.timeoutTime.Sub(time.Now())
	secondsLeft := timeDifference.Seconds()
	cappedSeconds := math.Max(0, secondsLeft)
	return int(cappedSeconds)
}

func (game *Game) addPlayer(player *Player) {
	players := make([]MessagePlayer, 0, len(game.Players))
	msgPlayer := MessagePlayer{
		Id:      player.TempId,
		Name:    player.Name,
		IsGuest: player.IsGuest,
		Score:   game.score[player],
	}

	joinMessage := NewPlayerJoinMessage(msgPlayer)
	for _, p := range game.Players {
		p.OutData <- joinMessage
		players = append(players, MessagePlayer{
			Id:      p.TempId,
			Name:    p.Name,
			IsGuest: p.IsGuest,
			Score:   game.score[p],
		})
	}

	activeState := game.ActiveState()

	player.OutData <- NewWelcomeMessage(msgPlayer, "")
	player.OutData <- NewGameStateMessage(activeState.Name(), activeState, players, game.timeLeft())
	game.Players = append(game.Players, player)
}

func (game *Game) RemovePlayer(player *Player) {
	game.Players = game.Players.Remove(player)
	game.Broadcast(NewPlayerLeaveMessage(player.Name))
}

func (game *Game) Strokes(from *Player, strokes []Stroke) {
	strokesMessage := NewStrokesMessage(strokes)
	for _, p := range game.Players {
		if p != from {
			p.OutData <- strokesMessage
		}
	}
}

func (game *Game) AddScore(player *Player, score int) {
	currentScore := game.score[player]
	game.score[player] = currentScore + score
}

func (game *Game) Run() {
	game.ActiveState().Activate(game)
	for len(game.state) > 0 && game.isRunning {
		select {
		case <-game.timeout:
			fmt.Println("Timeout occurred")
			game.ActiveState().Timeout()
		case player := <-game.playerJoin:
			game.addPlayer(player)
			game.ActiveState().PlayerJoin(player)
		case player := <-game.playerLeave:
			fmt.Println("Player left (game)")
			player.HasLeft = true
			game.ActiveState().PlayerLeave(player)
			game.RemovePlayer(player)
			if len(game.Players) == 0 {
				game.Stop()
			}
		case message, ok := <-game.inData:
			if !ok {
				break
			} else if !game.handleMessage(message) {
				game.ActiveState().Message(message.Source, message.Message)
			}
		}
	}
	for _, player := range game.Players {
		close(player.OutData)
	}
	close(game.closed)
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

func (game *Game) handleMessage(message inMessage) bool {
	err := game.messageHandler.Handle(message.Source, message.Message)
	return err == nil
}
