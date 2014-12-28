package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var FailedApiRequestError = errors.New("Request to server failed") //todo: change this error handling to actually use error sent by server

func parseResponse(response *http.Response, result interface{}, resultErr interface{}) error {
	defer response.Body.Close()
	if responseBody, err := ioutil.ReadAll(response.Body); err != nil {
		fmt.Println("Error while reading response body", err)
		return err
	} else {
		if response.StatusCode != http.StatusOK {
			if err := json.Unmarshal(responseBody, resultErr); err != nil {
				fmt.Println("Error while unmarshaling response body ", string(responseBody), err)
				return err
			}
		} else {
			if err := json.Unmarshal(responseBody, result); err != nil {
				fmt.Println("Error while unmarshaling response body ", string(responseBody), err)
				return err
			}
		}
	}
	return nil
}

func Post(address string, payload interface{}, result interface{}, resultErr interface{}) error {
	var requestBody []byte
	if payload != nil {
		var err error
		if requestBody, err = json.Marshal(payload); err != nil {
			fmt.Println("Error while marshaling request data ", err)
			return err
		}
	} else {
		requestBody = []byte{}
	}
	if response, err := http.Post(address, "text/json", bytes.NewReader(requestBody)); err != nil {
		fmt.Println("Error while sending request ", err)
		return err
	} else {
		if err := parseResponse(response, result, resultErr); err != nil {
			fmt.Println("Error while parsing reseponse ", err)
			return err
		}
	}
	return nil
}

func Get(address string, result interface{}, resultErr interface{}) error {
	if response, err := http.Get(address); err != nil {
		return err
	} else {
		if err := parseResponse(response, result, resultErr); err != nil {
			fmt.Println("Error while parsing response ", err)
			return err
		}
	}
	return nil
}
