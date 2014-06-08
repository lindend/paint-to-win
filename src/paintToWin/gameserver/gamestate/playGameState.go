package gamestate

import (
	"paintToWin/gameserver/game"
	"time"
)

const MinNumPlayers = 3

type PlayGameState struct {
	DefaultDeactivate
	DefaultPlayerJoin

	game    *game.Game
	context stateContext

	correctGuessers []*game.Player

	messageHandler *game.MessageHandler

	DrawingPlayerId string
}

func newPlayGameState(context stateContext) *PlayGameState {
	playState := &PlayGameState{
		context:         context,
		messageHandler:  game.NewMessageHandler(),
		DrawingPlayerId: context.drawingPlayer.TempId,
	}

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
}

func (p PlayGameState) Message(message game.InMessage) {
	p.messageHandler.Handle(message)
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
		p.game.Strokes(player, strokes.Strokes)
	}
}

func (p *PlayGameState) guessMessage(player *game.Player, guess *game.GuessMessage) {
	if player != p.context.drawingPlayer && player != p.context.choosingPlayer {
		if guess.Guess == p.context.word {
			p.game.Broadcast(game.NewCorrectGuessMessage(player.Name))
			if len(p.correctGuessers) == 0 {
				p.game.SetTimeout(10 * time.Second)
			}
			p.correctGuessers = append(p.correctGuessers, player)
		} else {
			p.game.Broadcast(game.NewWrongGuessMessage(player.Name, guess.Guess))
		}
	}
}
