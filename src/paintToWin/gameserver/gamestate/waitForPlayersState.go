package gamestate

import (
	"paintToWin/gameserver/game"
	"time"
)

type WaitForPlayersState struct {
	DefaultDeactivate
	DefaultMessageHandling

	game *game.Game

	MinNumPlayers int
	timeout       time.Duration

	context stateContext
}

func newWaitForPlayersState(minNumPlayers int, timeout time.Duration, context stateContext) *WaitForPlayersState {
	return &WaitForPlayersState{
		MinNumPlayers: minNumPlayers,
		timeout:       timeout,
		context:       context,
	}
}

func (w WaitForPlayersState) Timeout() {
	w.game.Stop()
}

func (w *WaitForPlayersState) Activate(game *game.Game) {
	w.game = game
	w.game.SetTimeout(w.timeout)
}

func (w WaitForPlayersState) PlayerJoin(player *game.Player) {
	if len(w.game.Players) >= w.MinNumPlayers {
		w.game.PopState()
	}
}

func (w WaitForPlayersState) PlayerLeave(player *game.Player) {
	if player == w.context.drawingPlayer {
		w.context.drawingPlayer = w.game.NextPlayer(w.context.drawingPlayer)
	}
}
