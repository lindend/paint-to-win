package storage

type SessionPlayer struct {
	PlayerName string
	IsGuest    bool
	UserId     int64
}

type Session struct {
	Id     string
	Player *SessionPlayer
}
