package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	"paintToWin/lobby/api"
	"paintToWin/server"
	"paintToWin/settings"
	"paintToWin/storage"
)

func main() {
	var dbConnectionString string

	flag.StringVar(&dbConnectionString, "db", "", "connection string for the database")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Initializing db")
	database, err := storage.InitializeDatabase(dbConnectionString)
	if err != nil {
		log.Fatal("Unable to initialize db ", err)
		return
	}

	serverInfo, err := server.LoadServerInfo()
	if err != nil {
		log.Fatal("Unable to load server info")
		return
	}

	fmt.Println("Loading config")
	config := Config{}
	if err = settings.Load(serverInfo.Name, database, &config); err != nil {
		log.Fatal("Error while loading config: \n" + err.Error())
		return
	}

	fmt.Println("")
	store, err := storage.NewStorage(&database, config.RedisAddress)
	if err != nil {
		log.Fatal("Unable to initialize storage ", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	api.RegisterUserApi(router, store)
	api.RegisterGameApi(router, store)

	fmt.Println("Listening on port ")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.LobbyApiPort), router))
}
