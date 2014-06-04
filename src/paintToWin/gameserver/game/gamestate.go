package game

type GameState interface {
	Name() string
	Timeout()
	Activate(game *Game)
	Deactivate()
	Message(message InMessage)
	PlayerJoin(player *Player)
	PlayerLeave(player *Player)
}
