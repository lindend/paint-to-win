package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"paintToWin/service"
	"paintToWin/storage"
	"paintToWin/web"
	"paintToWin/wordlistService/api"
)

func GetWordlistsHandler() web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		var output api.GetWordlistsOutput
		err := service.FindAndCall(api.GetWordlistsOperation, nil, &output)
		if err != nil {
			fmt.Println("Unable to call wordlist service", err)
		}
		return output, nil
	}
}

func RegisterMetadataApi(router *mux.Router, store *storage.Storage) {
	authenticator := NewSessionAuthenticator(store)

	router.HandleFunc("/wordlists", web.DefaultAuthenticateHandler(GetWordlistsHandler(), authenticator)).Methods("GET", "OPTIONS")
}
