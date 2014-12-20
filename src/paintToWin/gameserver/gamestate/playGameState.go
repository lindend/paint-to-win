package gamestate

import (
	"strings"
	"time"

	"paintToWin/gameserver/game"
)

const MinNumPlayers = 3

type PlayGameState struct {
	DefaultDeactivate

	game    *game.Game
	context stateContext

	messageHandler *game.MessageHandler

	DrawingPlayerId string

	roundStrokes []game.Stroke

	word string
}

func newPlayGameState(context stateContext) *PlayGameState {
	playState := &PlayGameState{
		context:         context,
		messageHandler:  game.NewMessageHandler(),
		DrawingPlayerId: context.drawingPlayer.TempId,
		word:            context.words[0],
	}
	context.words = context.words[1:]

	playState.messageHandler.Add(playState.guessMessage)
	playState.messageHandler.Add(playState.strokesMessage)

	return playState
}

func (p PlayGameState) Name() string {
	return "PlayGameState"
}

func (p *PlayGameState) Timeout() {
	p.game.SwapState(newEndRoundState(p.context))
}

func (p *PlayGameState) Activate(g *game.Game) {
	p.game = g
	p.game.SetTimeout(120 * time.Second)
	p.context.drawingPlayer.OutData <- game.NewTurnToPaintMessage(p.word)
}

func (p PlayGameState) Message(source *game.Player, message game.Message) {
	p.messageHandler.Handle(source, message)
}

func (p PlayGameState) PlayerJoin(player *game.Player) {
	allStrokesMessage := game.NewStrokesMessage(p.roundStrokes)
	player.OutData <- allStrokesMessage
}

func (p PlayGameState) PlayerLeave(player *game.Player) {
	if player == p.context.drawingPlayer {
		p.game.SwapState(newEndRoundState(p.context))
	} else if len(p.game.Players)-1 < MinNumPlayers {
		p.game.SwapState(newEndRoundState(p.context))
	}
}

func (p *PlayGameState) strokesMessage(player *game.Player, strokes *game.StrokesMessage) {
	if player == p.context.drawingPlayer {
		for _, stroke := range strokes.Strokes {
			p.roundStrokes = append(p.roundStrokes, stroke)
		}
		p.game.Strokes(player, strokes.Strokes)
	}
}

func (p *PlayGameState) guessMessage(player *game.Player, guess *game.GuessMessage) {
	if player == p.context.drawingPlayer {
		return
	}

	if p.context.correctGuessers.Contains(player) {
		return
	}

	if strings.ToLower(guess.Guess) == strings.ToLower(p.word) {
		p.game.Broadcast(game.NewCorrectGuessMessage(player.TempId))
		if len(p.context.correctGuessers) == 0 {
			p.game.SetTimeout(10 * time.Second)
		}
		p.context.correctGuessers = append(p.context.correctGuessers, player)
	} else {
		p.game.Broadcast(game.NewWrongGuessMessage(player.TempId, guess.Guess))
	}

}
