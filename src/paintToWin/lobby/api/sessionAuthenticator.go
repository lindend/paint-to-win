package api

import (
	"errors"
	"fmt"
	"net/http"
	"paintToWin/storage"

	"github.com/gorilla/context"
)

type SessionAuthenticator struct {
	store *storage.Storage
}

func NewSessionAuthenticator(store *storage.Storage) *SessionAuthenticator {
	return &SessionAuthenticator{store}
}

const SessionAuthorization = "Session"

var UnknownAuthorizationType = errors.New("Unrecognized authorization header")
var InvalidCredentials = errors.New("Invalid credentials")

func (auth SessionAuthenticator) Authenticate(authorizationType string, authorization string, req *http.Request) error {
	fmt.Println("Authorizing ", authorizationType, authorization)
	if authorizationType == SessionAuthorization {
		if len(authorization) == 0 {
			return InvalidCredentials
		}
		session := &storage.Session{}
		if err := auth.store.GetFromCache("session:"+authorization, session); err != nil {
			return err
		}
		context.Set(req, "session", session)
		return nil
	} else {
		return UnknownAuthorizationType
	}
}
