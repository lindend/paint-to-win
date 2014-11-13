package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"paintToWin/gameserver/api"
	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/gameserver/network/ws"
	"paintToWin/server"
	"paintToWin/settings"
	"paintToWin/storage"
)

func main() {
	var dbConnectionString string
	var cpuprofile string

	flag.StringVar(&dbConnectionString, "db", "", "connection string for the database")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "path to cpu profile output")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	idGenerator := StartIdGenerator()

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

	fmt.Println("Server info: ", serverInfo)

	config := Config{}
	if err = settings.Load(serverInfo.Name, &database, &config); err != nil {
		log.Fatal("Error while loading config: \n" + err.Error())
		return
	}

	currentServer := storage.Server{
		Name:    serverInfo.Name,
		Address: fmt.Sprintf("http://%v:%d", serverInfo.HostName, config.GameServerApiPort),
		Type:    "gameserver",
	}

	database.Where(storage.Server{Name: currentServer.Name}).Assign(currentServer).FirstOrInit(&currentServer)
	database.Save(&currentServer)

	store, err := storage.NewStorage(&database, config.RedisAddress)
	if err != nil {
		log.Fatal("Unable to initialize storage ", err)
		return
	}

	connectChan := make(chan network.NewConnection)

	endpoints := []network.EndpointInfo{}

	err = startWebsocketEndpoint(config.GameServerGamePort, connectChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	endpoints = append(endpoints, network.EndpointInfo{
		Address:  config.Address,
		Port:     config.GameServerGamePort,
		Protocol: "ws",
	})

	gameManager := gamemanager.NewGameManager(idGenerator, endpoints, store, currentServer)

	go func() {
		for connection := range connectChan {
			ClientHandshake(gameManager, store, idGenerator, connection)
		}
	}()

	if err := api.Start(config.GameServerApiPort, gameManager); err != nil {
		log.Fatal("Error while initializing web API ", err)
		return
	}

	fmt.Println("Initializing profiling")
	if err := initProfiling(cpuprofile); err != nil {
		fmt.Println("Could not initialize profiling", err)
	} else {
		fmt.Println("Profiling initialized")
		defer func() {
			fmt.Println("Stopping profiling")
			stopProfiling()
		}()
	}

	fmt.Println("Game server successfully started")
	//for {
	time.Sleep(1000)
	ioutil.ReadAll(os.Stdin)
	//}
}

func startWebsocketEndpoint(port int, onConnect chan network.NewConnection) error {
	fmt.Println("Starting web socket server")
	wsEndpoint, err := ws.StartWebSocketServer(port, []string{"/{reservationId}/{sessionId}"})
	if err != nil {
		return err
	}

	mergeNewConnection(wsEndpoint.OnConnect, onConnect)
	fmt.Println("Successfully started web socket server")
	return nil
}

func mergeNewConnection(newConnect <-chan network.NewConnection, merged chan<- network.NewConnection) {
	go func() {
		for {
			conn := <-newConnect
			merged <- conn
		}
	}()
}
