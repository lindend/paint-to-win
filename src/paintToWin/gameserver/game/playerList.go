package game

type PlayerList []*Player

func (pl PlayerList) Contains(player *Player) bool {
	_, err := pl.FindPlayer(player)
	return err == nil
}

func (pl PlayerList) FindPlayer(player *Player) (int, error) {
	for i, p := range pl {
		if p == player {
			return i, nil
		}
	}
	return 0, PlayerNotFoundError
}

func (pl PlayerList) NextPlayer(player *Player) *Player {
	playerIdx, err := pl.FindPlayer(player)
	if err != nil {
		return nil
	}
	for i := playerIdx + 1; i < len(pl); i++ {
		if !pl[i].HasLeft {
			return pl[i]
		}
	}

	for i := 0; i < playerIdx; i++ {
		if !pl[i].HasLeft {
			return pl[i]
		}
	}
	return nil
}

func (pl PlayerList) PreviousPlayer(player *Player) *Player {
	playerIdx, err := pl.FindPlayer(player)
	if err != nil {
		return nil
	}
	for i := playerIdx - 1; i > 0; i-- {
		if !pl[i].HasLeft {
			return pl[i]
		}
	}

	for i := len(pl) - 1; i > playerIdx; i-- {
		if !pl[i].HasLeft {
			return pl[i]
		}
	}
	return nil
}

func (pl PlayerList) Remove(player *Player) PlayerList {
	i, err := pl.FindPlayer(player)
	if err != nil {
		return pl
	}
	return append(pl[:i], pl[i+1:]...)
}
