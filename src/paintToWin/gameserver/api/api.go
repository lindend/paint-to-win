package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"paintToWin/gameserver/gamemanager"
	"paintToWin/service"
)

const serviceName = "gameserver"
const Service = serviceName

func Start(location service.Location, gameManager *gamemanager.GameManager) error {
	router := mux.NewRouter()

	host := service.NewHttpServiceHost(location, router)
	RegisterGameManagerApi(host, gameManager)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", location.Port), router))
	}()

	return nil
}
