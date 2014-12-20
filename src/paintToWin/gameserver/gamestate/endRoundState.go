package gamestate

import "paintToWin/gameserver/game"

type EndRoundState struct {
	DefaultTimeout
	DefaultDeactivate
	DefaultMessageHandling
	DefaultPlayerJoin
	DefaultPlayerLeave

	game    *game.Game
	context stateContext
}

func newEndRoundState(context stateContext) *EndRoundState {
	return &EndRoundState{
		context: context,
	}
}

func (e EndRoundState) Name() string {
	return "EndRoundState"
}

func (e *EndRoundState) Activate(g *game.Game) {
	e.game = g

	e.game.AddScore(e.context.drawingPlayer, calculateDrawingPlayerScore(e.context.correctGuessers, e.game.Players))

	for _, pl := range e.context.correctGuessers {
		e.game.AddScore(pl, calculateCorrectPlayerScore(e.context.correctGuessers, e.game.Players))
	}

	if e.game.CurrentRound < e.game.NumRounds {
		e.context.drawingPlayer = e.game.Players.NextPlayer(e.context.drawingPlayer)
		e.game.SwapState(newInitRoundState(e.context))
	} else {
		e.game.Stop()
	}
}

func calculateDrawingPlayerScore(correctGuessers []*game.Player, players []*game.Player) int {
	numCorrect := len(correctGuessers)
	numPlayers := len(players)

	if numCorrect == 0 {
		return 0
	}

	if numCorrect > numPlayers/3 {
		return 2
	} else if numCorrect > (2*numPlayers)/3 {
		return 3
	}
	return 1
}

func calculateCorrectPlayerScore(correctGuessers, players []*game.Player) int {
	return 2
}
