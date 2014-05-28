package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paintToWin/gameserver/gamemanager"
)

func Start(apiPort int, gameManager *gamemanager.GameManager) error {
	router := mux.NewRouter()

	RegisterGameManagerApi(router, gameManager)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), router))
	}()

	return nil
}
