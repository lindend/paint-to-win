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

	"paintToWin/server"
	"paintToWin/service"
	"paintToWin/settings"
	"paintToWin/storage"
	"paintToWin/wordlist"
	"paintToWin/wordlistService/api"
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

	currentServer := storage.Server{
		Name:    serverInfo.Name,
		Address: fmt.Sprintf("http://%v:%d", config.Address, config.ApiPort),
		Type:    "wordlist",
	}

	database.Where(storage.Server{Name: currentServer.Name}).Assign(currentServer).FirstOrInit(&currentServer)
	database.Save(&currentServer)

	location := service.Location{
		Address: config.Address,
		Port:    config.ApiPort,

		Protocol:  "HTTP",
		Transport: "TCP",

		Priority: 0,
		Weight:   0,
	}

	serviceManager := service.NewDbServiceManager(&database)
	serviceManager.Register(api.Service, location)

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	wordlistInfos, err := wordlist.EnumerateWordlists(config.WordlistRoot)
	if err != nil {
		log.Fatal("Unable to find wordlists ", err)
	}

	wordlists := make([]wordlist.Wordlist, 0)
	for _, wordlistInfo := range wordlistInfos {
		fmt.Println("Loading wordlist from", wordlistInfo.Path)
		if wordlist, err := wordlist.LoadWordlistFromFile(wordlistInfo); err != nil {
			fmt.Println("Unable to load wordlist from " + wordlistInfo.Path)
		} else {
			wordlists = append(wordlists, wordlist)
		}
	}

	host := service.NewHttpServiceHost(location, router)

	fmt.Println("Initializing API")
	api.RegisterWordlistApi(host, wordlists)

	fmt.Sprintln("Listening on port %v", config.ApiPort)
	fmt.Println("Starting web service")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.ApiPort), router))
}
