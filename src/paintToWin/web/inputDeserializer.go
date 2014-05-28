package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type InputError struct {
	field   string
	message string
}

func (inputError InputError) Error() string {
	return inputError.field + ": " + inputError.message
}

func NewInputError(field string, message string) InputError {
	return InputError{
		field:   field,
		message: message,
	}
}

type Validater interface {
	Validate() []InputError
}

var UnknownContentTypeError = errors.New("Unable to find a decoder for specified content type")

func DeserializeInput(req *http.Request, value interface{}) error {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/json") || strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(req.Body)
		return decoder.Decode(value)
	}
	return UnknownContentTypeError
}

func DeserializeAndValidateInput(req *http.Request, value Validater) ([]InputError, error) {
	if err := DeserializeInput(req, value); err != nil {
		return nil, err
	} else {
		return value.Validate(), nil
	}
}
