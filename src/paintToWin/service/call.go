package service

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"paintToWin/web"
)

func FindAndCall(operation ServiceOperation, input interface{}) (interface{}, error) {
	locations, err := Find(operation.ServiceName)
	if err != nil {
		return nil, err
	}
	return Call(operation, locations[0], input)
}

func Call(operation ServiceOperation, location Location, input interface{}) (interface{}, error) {
	if !verifyType(operation.InputType, input) {
		panic("Invalid input type")
	}

	output := reflect.Zero(operation.OutputType).Interface()

	switch location.Protocol {
	case HttpProtocol:
		err := callHttpService(operation, location, input, output)
		return output, err
	}
	return nil, errors.New("No such service type available")
}

func callHttpService(operation ServiceOperation, location Location, input interface{}, output interface{}) error {
	path := resolveOperationPath(operation.Path, input)
	operationUri := fmt.Sprintf("%v://%v:%v/%v", location.Protocol, location.Address, location.Port, path)

	switch operation.Method {
	case "POST":
		var errResult ServiceError
		return web.Post(operationUri, input, output, &errResult)
	case "GET":
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
