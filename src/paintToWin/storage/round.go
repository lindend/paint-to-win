package storage

type Round struct {
	Id     int64
	Index  int
	GameId int64 `sql:not null`
	Word   string
}
