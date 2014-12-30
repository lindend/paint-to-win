package service

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"paintToWin/web"
)

func FindAndCall(operation ServiceOperation, input interface{}, output interface{}) error {
	locations, err := Find(operation.ServiceName)
	if err != nil {
		return err
	}
	return Call(operation, locations[0], input, output)
}

func Call(operation ServiceOperation, location Location, input interface{}, output interface{}) error {
	if !verifyType(operation.InputType, input) {
		panic("Invalid input type")
	}

	switch strings.ToLower(location.Protocol) {
	case HttpProtocol:
		err := callHttpService(operation, location, input, output)
		return err
	}
	return errors.New("Protocol not supported: " + location.Protocol)
}

func callHttpService(operation ServiceOperation, location Location, input interface{}, output interface{}) error {
	path := resolveOperationPath(operation.Path, input)
	operationUri := fmt.Sprintf("%v://%v:%v/%v", location.Protocol, location.Address, location.Port, path)

	switch strings.ToLower(operation.Method) {
	case "post":
		var errResult ServiceError
		return web.Post(operationUri, input, output, &errResult)
	case "get":
		var errResult ServiceError
		return web.Get(operationUri, output, &errResult)
	}
	panic("Invalid operation type for HTTP service: " + operation.Method)
}

func verifyType(t reflect.Type, obj interface{}) bool {
	objType := reflect.TypeOf(obj)
	return objType == t
}

func resolveOperationPath(operation string, input interface{}) string {
	inputValue := reflect.ValueOf(input)

	re := regexp.MustCompile("{.+?}")
	return re.ReplaceAllStringFunc(operation, func(match string) string {
		fieldName := match[1 : len(match)-1]
		field := inputValue.FieldByName(fieldName)
		if !field.IsValid() {
			panic("Invalid field " + fieldName + " in operation " + operation)
		}

		fieldValue := fmt.Sprintf("%v", field.Interface())

		return fieldValue
	})
}
