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

func (e *EndRoundState) Activate(g *game.Game) {
	e.game = g

	if e.game.CurrentRound < e.game.NumRounds {
		e.context.drawingPlayer = e.game.NextPlayer(e.context.drawingPlayer)
		e.game.SwapState(newInitRoundState(e.context))
	} else {
		e.game.Stop()
	}
}
