package storage

type Wordlist struct {
	Id       int64
	Name     string `sql:not null`
	Language string `sql:not null`
}
