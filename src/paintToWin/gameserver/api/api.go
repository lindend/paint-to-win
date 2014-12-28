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

func Start(address string, apiPort int, gameManager *gamemanager.GameManager, serviceManager service.ServiceManager) error {
	router := mux.NewRouter()

	location := service.Location{
		Address: address,
		Port:    apiPort,

		Protocol:  "HTTP",
		Transport: "TCP",

		Priority: 0,
		Weight:   0,
	}

	host := service.NewHttpServiceHost(location, serviceManager, router)
	RegisterGameManagerApi(host, gameManager)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", location.Port), router))
	}()

	return nil
}
