package storage

type Game struct {
	Id         int64
	GameId     string `sql:not null;unique`
	Name       string `sql:not null`
	CreatedBy  string `sql:not null`
	HostedOn   string `sql:not null`
	IsActive   bool   `sql:not null`
	NumPlayers int    `sql:"-"`
}
