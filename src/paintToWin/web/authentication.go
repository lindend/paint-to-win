package web

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator interface {
	Authenticate(authType string, authentication string, req *http.Request) error
}

func Authenticate(handler http.HandlerFunc, authenticator Authenticator) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		fmt.Println("Authenticating ", authorizationHeader)
		if authorizationHeader != "" {
			splitAuthorization := strings.SplitN(authorizationHeader, " ", 2)
			if len(splitAuthorization) != 2 {
				rw.WriteHeader(http.StatusUnauthorized)
			} else if err := authenticator.Authenticate(splitAuthorization[0], splitAuthorization[1], req); err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
			} else {
				handler(rw, req)
			}
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
		}
	}
}
