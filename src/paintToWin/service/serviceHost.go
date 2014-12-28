package service

type Host interface {
	Register(function interface{}, operation ServiceOperation) error
}
