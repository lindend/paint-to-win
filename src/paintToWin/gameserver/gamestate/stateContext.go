package gamestate

import "paintToWin/gameserver/game"

type stateContext struct {
	drawingPlayer   *game.Player
	choosingPlayer  *game.Player
	word            string
	correctGuessers game.PlayerList
}
