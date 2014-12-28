package storage

type Service struct {
	Id   int64
	Name string

	Address string
	Port    int

	Protocol  string
	Transport string

	Priority int
	Weight   int
}
