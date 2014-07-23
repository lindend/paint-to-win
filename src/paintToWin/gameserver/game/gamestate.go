package game

type GameState interface {
	Name() string
	Timeout()
	Activate(game *Game)
	Deactivate()
	Message(player *Player, message Message)
	PlayerJoin(player *Player)
	PlayerLeave(player *Player)
}
