package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiError interface{
	Error() interface{}
	HttpStatusCode() int
}

type apiError struct {
	error interface{}
	httpStatusCode int
}

func (err apiError) Error() interface{} {
	return err.error
}

func (err apiError) HttpStatusCode() int {
	return err.httpStatusCode
}

func NewApiError(statusCode int, err interface{}) ApiError {
	return apiError {
		error: err,
		httpStatusCode: statusCode,
	}
}

type RequestHandler func(request *http.Request) (interface{}, ApiError)


func OutputSerializer(handler RequestHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		data, err := handler(req)

		if err != nil {
			fmt.Println("Http request failed", err.HttpStatusCode(), err)

			rw.WriteHeader(err.HttpStatusCode())
			if jsonData, serializeErr := json.Marshal(err.Error()); serializeErr != nil {
				rw.Write([]byte("Error while serializing response"))
			} else {
				rw.Write(jsonData)
			}
		} else {
			if jsonData, serializeErr := json.Marshal(data); serializeErr != nil {
				rw.WriteHeader(500)
				rw.Write([]byte("Error while serializing response"))
			} else {
				rw.WriteHeader(http.StatusOK)
				rw.Write(jsonData)
			}
		}


	}
}
