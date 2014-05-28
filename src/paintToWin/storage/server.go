package storage

import "time"

type Server struct {
	Id      int64
	Name    string `sql:not null;unique`
	Address string `sql:not null`
	Type    string `sql:not null` //gameserver or lobbyserver

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
