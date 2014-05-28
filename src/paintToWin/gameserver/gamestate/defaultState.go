package gamestate

import "paintToWin/gameserver/game"

type DefaultTimeout struct {
}

func (d DefaultTimeout) Timeout() {
}

type DefaultActivate struct {
	game *game.Game
}

func (d *DefaultActivate) Activate(game *game.Game) {
	d.game = game
}

type DefaultDeactivate struct {
}

func (d DefaultDeactivate) Deactivate() {
}

type DefaultMessageHandling struct {
}

func (d DefaultMessageHandling) Message(message game.InMessage) {
}

type DefaultPlayerJoin struct {
}

func (d DefaultPlayerJoin) PlayerJoin(player *game.Player) {
}

type DefaultPlayerLeave struct {
}

func (d DefaultPlayerLeave) PlayerLeave(player *game.Player) {
}
