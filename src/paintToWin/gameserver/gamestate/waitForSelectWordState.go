package gamestate

import (
	"paintToWin/gameserver/game"
)

type WaitForSelectWordState struct {
	DefaultActivate
	DefaultDeactivate
	DefaultPlayerJoin

	context stateContext

	messageHandler *game.MessageHandler

	SelectingPlayer string
}

func newWaitForSelectWordState(context stateContext) *WaitForSelectWordState {
	waitState := &WaitForSelectWordState{
		context:         context,
		messageHandler:  game.NewMessageHandler(),
		SelectingPlayer: context.choosingPlayer.TempId,
	}
	waitState.messageHandler.Add(waitState.chooseWordMessage)
	return waitState
}

func (w *WaitForSelectWordState) chooseWordMessage(player *game.Player, message *game.ChooseWordMessage) {
	if player == w.context.choosingPlayer {
		w.context.word = message.Word
		w.game.SwapState(newPlayGameState(w.context))
		return
	}
}

func (w WaitForSelectWordState) Name() string {
	return "WaitForWordState"
}

func (w WaitForSelectWordState) Timeout() {
	w.game.Stop()
}

func (w WaitForSelectWordState) Message(message game.InMessage) {
	w.messageHandler.Handle(message)
}

func (w WaitForSelectWordState) PlayerLeave(player *game.Player) {
	if player == w.context.drawingPlayer {
		w.context.drawingPlayer = w.game.NextPlayer(w.context.drawingPlayer)
		w.game.SwapState(newInitRoundState(w.context))
	} else if player == w.context.choosingPlayer {
		w.game.SwapState(newInitRoundState(w.context))
	} else if len(w.game.Players) < 3 {
		w.game.SwapState(newInitRoundState(w.context))
	}
}
