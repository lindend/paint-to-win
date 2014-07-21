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
	"paintToWin/gameserver/codec"
	"paintToWin/gameserver/communication"
	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/gameserver/network/ws"
	"paintToWin/server"
	"paintToWin/settings"
	"paintToWin/storage"
)

func main() {
	var dbConnectionString string
	var address string
	var cpuprofile string

	flag.StringVar(&dbConnectionString, "db", "", "connection string for the database")
	flag.StringVar(&address, "address", "", "remotely accessible address of the server")
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
	if err = settings.Load(serverInfo.Name, database, &config); err != nil {
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

	messageChan := make(chan communication.Message)
	connectChan := make(chan network.NewConnection)
	disconnectChan := make(chan network.Connection)

	endpoints := []network.EndpointInfo{}

	err = startWebsocketEndpoint(config.GameServerGamePort, messageChan, connectChan, disconnectChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	endpoints = append(endpoints, network.EndpointInfo{
		Address:  address,
		Port:     config.GameServerGamePort,
		Protocol: "ws",
	})

	commHub, commOutData := startCommunicationHub(messageChan, connectChan, disconnectChan)
	gameManager := gamemanager.NewGameManager(idGenerator, endpoints, commHub, store, currentServer)

	handshake := CreateClientHandshake(gameManager, store, idGenerator)

	fmt.Println("Starting comm hub serve")
	go commHub.Serve(handshake)

	if err := api.Start(config.GameServerApiPort, gameManager); err != nil {
		log.Fatal("Error while initializing web API ", err)
		return
	}

	go func() {
		for {
			msg := <-commOutData
			fmt.Println("Packet output main.go")
			msg.Connection.Send(msg.Data)
		}
	}()

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

func startWebsocketEndpoint(port int, onMessage chan communication.Message, onConnect chan network.NewConnection, onDisconnect chan network.Connection) error {
	fmt.Println("Starting web socket server")
	wsEndpoint, err := ws.StartWebSocketServer(port, []string{"/{reservationId}/{sessionId}"})
	if err != nil {
		return err
	}
	wsCodec := codec.StandardDecoder(wsEndpoint.OnData)

	mergeOnMessage(wsCodec, onMessage)
	mergeNewConnection(wsEndpoint.OnConnect, onConnect)
	mergeConnect(wsEndpoint.OnDisconnect, onDisconnect)
	fmt.Println("Successfully started web socket server")
	return nil
}

func startCommunicationHub(
	onMessage chan communication.Message,
	onConnect chan network.NewConnection,
	onDisconnect chan network.Connection,
) (*communication.CommunicationHub, <-chan communication.Message) {

	fmt.Println("Starting communication hub")
	commHub, commChan := communication.NewCommunicationHub(onConnect, onDisconnect, onMessage)
	fmt.Println("Successfully started communication hub")
	return commHub, commChan
}

func mergeOnMessage(inNewMessage <-chan communication.Message, mergedChan chan<- communication.Message) {
	go func() {
		for {
			msg := <-inNewMessage
			fmt.Println("new message main.go")
			mergedChan <- msg
		}
	}()
}

func mergeNewConnection(newConnect <-chan network.NewConnection, merged chan<- network.NewConnection) {
	go func() {
		for {
			conn := <-newConnect
			merged <- conn
		}
	}()
}

func mergeConnect(newDisconnec <-chan network.Connection, mergedChan chan<- network.Connection) {
	go func() {
		for {
			conn := <-newDisconnec
			mergedChan <- conn
		}
	}()
}
