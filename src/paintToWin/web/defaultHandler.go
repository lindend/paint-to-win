package web

import (
	"net/http"
)

func DefaultHandler(requestHandler RequestHandler) http.HandlerFunc {
	return EnableCors(OutputSerializer(requestHandler))
}

func DefaultAuthenticateHandler(requestHandler RequestHandler, authenticator Authenticator) http.HandlerFunc {
	return EnableCors(Authenticate(OutputSerializer(requestHandler), authenticator))
}
