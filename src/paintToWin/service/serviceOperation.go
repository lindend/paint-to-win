package service

import (
	"reflect"
)

type ServiceOperation struct {
	ServiceName string
	Path        string
	Name        string
	Method      string

	InputType  reflect.Type
	OutputType reflect.Type
}

func NewOperation(serviceName string, path string, name string, method string, input interface{}, output interface{}) ServiceOperation {
	return ServiceOperation{
		ServiceName: serviceName,

		Name:   name,
		Path:   path,
		Method: method,

		InputType:  reflect.TypeOf(input),
		OutputType: reflect.TypeOf(output),
	}
}
