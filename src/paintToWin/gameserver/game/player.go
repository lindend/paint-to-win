package game

type Player struct {
	Name      string
	IsGuest   bool
	PlayerId  string
	SessionId string
	HasLeft   bool

	OutData chan<- Message
}

func NewPlayer(name string, isGuest bool, id string, sessionId string, outData chan<- Message) *Player {
	return &Player{
		Name:      name,
		IsGuest:   isGuest,
		PlayerId:  id,
		SessionId: sessionId,
		OutData:   outData,
		HasLeft:   false,
	}
}
