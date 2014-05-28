package web

import (
	"net/http"
)

func EnableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
		} else {
			handler(rw, req)
		}
	}
}
