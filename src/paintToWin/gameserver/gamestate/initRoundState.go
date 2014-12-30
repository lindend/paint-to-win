package gamestate

import (
	"paintToWin/gameserver/game"
	"time"
)

type InitRoundState struct {
	DefaultTimeout
	DefaultDeactivate
	DefaultPlayerJoin
	DefaultPlayerLeave
	DefaultMessageHandling
	game    *game.Game
	context stateContext
}

func NewInitRoundState(words []string) *InitRoundState {
	return &InitRoundState{
		context: stateContext{
			words: words,
		},
	}
}

func newInitRoundState(context stateContext) *InitRoundState {
	return &InitRoundState{
		context: context,
	}
}

func (state InitRoundState) Name() string {
	return "InitRoundState"
}

func (state *InitRoundState) Activate(g *game.Game) {
	state.game = g

	state.context.correctGuessers = nil

	if len(g.Players) < 3 {
		g.PushState(newWaitForPlayersState(MinNumPlayers, 100*time.Minute, state.context))
	} else {
		if state.context.drawingPlayer == nil {
			state.context.drawingPlayer = g.Players[0]
		}
		g.SwapState(newPlayGameState(state.context))
	}
}
