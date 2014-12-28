package util

import (
	"errors"
	"reflect"
	"strconv"
)

func ParseStringToValue(str string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.Bool:
		if value, err := strconv.ParseBool(str); err != nil {
			return errors.New("Cannot convert " + str + " to bool")
		} else {
			field.SetBool(value)
		}
	case reflect.Int:
		if value, err := strconv.ParseInt(str, 10, 32); err != nil {
			return errors.New("Cannot convert " + str + " to int")
		} else {
			field.SetInt(value)
		}
	case reflect.Int64:
		if value, err := strconv.ParseInt(str, 10, 64); err != nil {
			return errors.New("Cannot convert " + str + " to int64")
		} else {
			field.SetInt(value)
		}
	case reflect.String:
		field.SetString(str)
	default:
		return errors.New("No converter available for kind " + field.Kind().String())
	}
	return nil
}
