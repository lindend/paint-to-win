package service

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"paintToWin/util"
	"paintToWin/web"
)

type HttpServiceHost struct {
	serviceManager ServiceManager
	location       Location

	router *mux.Router
}

type inputBuilder func(req *http.Request) (reflect.Value, error)

func NewHttpServiceHost(location Location, serviceManager ServiceManager, router *mux.Router) *HttpServiceHost {
	return &HttpServiceHost{
		serviceManager: serviceManager,
		location:       location,
		router:         router,
	}
}

func (h *HttpServiceHost) Register(function interface{}, operation ServiceOperation) error {
	h.router.HandleFunc(operation.Path, web.DefaultHandler(buildHttpOperationHandler(function, operation, h.location))).Methods(operation.Method, "OPTIONS")
	return nil
}

func createInputBuilder(argument reflect.Type, operationInputs []string) inputBuilder {
	if argument == reflect.TypeOf(&http.Request{}) {
		return func(req *http.Request) (reflect.Value, error) {
			return reflect.ValueOf(req), nil
		}
	}

	if argument.Kind() != reflect.Struct {
		panic("Only structs are supported")
	}

	fields := make([]string, 0)

	for _, operationInput := range operationInputs {
		_, ok := argument.FieldByName(operationInput)

		if ok {
			fields = append(fields, operationInput)
		}
	}

	return func(req *http.Request) (reflect.Value, error) {
		inputValue := reflect.New(argument)
		vars := mux.Vars(req)

		web.DeserializeInput(req, inputValue.Interface())

		for _, field := range fields {
			varValue := vars[field]
			fieldValue := inputValue.FieldByName(field)

			err := util.ParseStringToValue(varValue, fieldValue)
			if err != nil {
				return inputValue, err
			}
		}
		return inputValue, nil
	}
}

func getHttpOperationInputs(operation string) []string {
	re := regexp.MustCompile("{.*?}")
	matches := re.FindAllString(operation, -1)
	result := make([]string, len(matches))

	for i, v := range matches {
		result[i] = strings.Trim(v, "{}")
	}
	return result
}

func parseHttpResults(values []reflect.Value) (interface{}, web.ApiError) {
	var errValue web.ApiError
	var resultValue interface{}

	for _, v := range values {
		if v.IsNil() {
			continue
		}

		if v.Type().Implements(reflect.TypeOf((*error)(nil))) {
			err := v.Interface().(error)
			errValue = web.NewApiError(http.StatusInternalServerError, err)
		} else {
			resultValue = v.Interface()
		}
	}

	return resultValue, errValue
}

func buildHttpOperationHandler(function interface{}, operation ServiceOperation, location Location) web.RequestHandler {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		panic("Http service function is not a function")
	}

	funcType := funcValue.Type()

	numInputs := funcType.NumIn()
	inputBuilders := make([]inputBuilder, numInputs)

	operationInputs := getHttpOperationInputs(operation.Path)

	for i := 0; i < numInputs; i++ {
		inputBuilders[i] = createInputBuilder(funcType.In(i), operationInputs)
	}

	return func(req *http.Request) (interface{}, web.ApiError) {
		inputs := make([]reflect.Value, numInputs)

		for i, v := range inputBuilders {
			value, err := v(req)
			if err != nil {
				return nil, web.NewApiError(http.StatusBadRequest, "")
			}
			inputs[i] = value
		}
		results := funcValue.Call(inputs)

		return parseHttpResults(results)
	}
}
