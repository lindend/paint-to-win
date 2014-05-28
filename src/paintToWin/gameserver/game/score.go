package game

type Score struct {
	Player *Player
	Score  int
}

func NewScore(player *Player, score int) Score {
	return Score{player, score}
}
