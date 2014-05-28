package main

import (
	"fmt"
	"log"
	"time"

	"paintToWin/gameserver/api"
	"paintToWin/gameserver/codec"
	"paintToWin/gameserver/communication"
	"paintToWin/gameserver/gamemanager"
	"paintToWin/gameserver/network"
	"paintToWin/gameserver/network/ws"
	"paintToWin/storage"
)

func main() {
	config := Config{
		WebsocketPort:      8085,
		ApiPort:            8084,
		DbConnectionString: "user=p2wuser password=devpassword host=xena port=5432 dbname=paint2win sslmode=disable",
		RedisAddress:       "10.10.0.98:6379",
	}
	idGenerator := StartIdGenerator()

	database, err := storage.InitializeDatabase(config.DbConnectionString)
	if err != nil {
		log.Fatal("Unable to initialize db ", err)
		return
	}

	var currentServer storage.Server
	var serverAddress string
	currentServer, serverAddress, err = loadServerInfo(config.ApiPort)
	if err != nil {
		log.Fatal("Error while reading server info", err)
		return
	}
	database.Where(storage.Server{Name: currentServer.Name}).Assign(currentServer).FirstOrInit(&currentServer)
	database.Save(&currentServer)

	store, err := storage.NewStorage(config.DbConnectionString, config.RedisAddress)
	if err != nil {
		log.Fatal("Unable to initialize storage ", err)
		return
	}

	messageChan := make(chan communication.Message)
	connectChan := make(chan network.NewConnection)
	disconnectChan := make(chan network.Connection)

	endpoints := []network.EndpointInfo{}

	err = startWebsocketEndpoint(config.WebsocketPort, messageChan, connectChan, disconnectChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	endpoints = append(endpoints, network.EndpointInfo{
		Address:  serverAddress,
		Port:     config.WebsocketPort,
		Protocol: "ws",
	})

	commHub, commOutData := startCommunicationHub(messageChan, connectChan, disconnectChan)
	gameManager := gamemanager.NewGameManager(idGenerator, endpoints, commHub, store, currentServer)

	handshake := CreateClientHandshake(gameManager, store, idGenerator)

	fmt.Println("Starting comm hub serve")
	go commHub.Serve(handshake)

	if err := api.Start(config.ApiPort, gameManager); err != nil {
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

	fmt.Println("Game server successfully started")
	for {
		time.Sleep(1000)
	}
}

func startWebsocketEndpoint(port int, onMessage chan communication.Message, onConnect chan network.NewConnection, onDisconnect chan network.Connection) error {
	fmt.Println("Starting web socket server")
	wsEndpoint, err := ws.StartWebSocketServer(port, []string{"/{reservationId}/{playerId}"})
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
