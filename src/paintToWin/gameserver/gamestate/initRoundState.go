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

func NewInitRoundState() *InitRoundState {
	return &InitRoundState{
		context: stateContext{},
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

	state.context.word = ""

	if len(g.Players) < 3 {
		g.PushState(newWaitForPlayersState(MinNumPlayers, 100*time.Minute, state.context))
	} else {
		if state.context.drawingPlayer == nil {
			state.context.drawingPlayer = g.Players[0]
		}
		state.context.choosingPlayer = g.PreviousPlayer(state.context.drawingPlayer)

		g.SwapState(newWaitForSelectWordState(state.context))
	}
}
