package gamestate

import "paintToWin/gameserver/game"

type stateContext struct {
	drawingPlayer   *game.Player
	words           []string
	correctGuessers game.PlayerList
}
