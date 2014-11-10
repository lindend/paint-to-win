package storage

type WordList struct {
	Id       int64
	Name     string `sql:not null`
	Language string `sql:not null`
}
