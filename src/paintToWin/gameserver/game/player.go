package game

type Player struct {
	Name     string
	IsGuest  bool
	PlayerId string
	TempId   string
	HasLeft  bool

	OutData chan<- Message
}

func NewPlayer(name string, isGuest bool, id string, tempId string, outData chan<- Message) *Player {
	return &Player{
		Name:     name,
		IsGuest:  isGuest,
		PlayerId: id,
		TempId:   tempId,
		OutData:  outData,
		HasLeft:  false,
	}
}
