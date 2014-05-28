package storage

import "time"

type Player struct {
	Id           int64
	UserName     string `sql:not null;unique`
	Email        string `sql: not null`
	PasswordHash []byte `sql: not null`
	Salt         []byte `sql: not null`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
