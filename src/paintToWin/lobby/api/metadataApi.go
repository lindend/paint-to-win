package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"paintToWin/lobby/game"
	"paintToWin/storage"
	"paintToWin/web"
)

func GetWordlistsHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		wordLists, err := game.GetWordLists(store)
		if err != nil {
			return nil, web.NewApiError(http.StatusInternalServerError, err.Error())
		}
		return wordLists, nil
	}
}

func RegisterMetadataApi(router *mux.Router, store *storage.Storage) {
	authenticator := NewSessionAuthenticator(store)

	router.HandleFunc("/wordlists", web.DefaultAuthenticateHandler(GetWordlistsHandler(store), authenticator)).Methods("GET", "OPTIONS")
}
