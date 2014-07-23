package ws

import (
	"fmt"
	"log"
	"net/http"
	"paintToWin/gameserver/network"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WsEndpoint struct {
	OnConnect chan network.NewConnection

	connections map[WsConnection]struct{}
}

type WsConnection struct {
	socket     *websocket.Conn
	connection network.Connection
}

func (endpoint *WsEndpoint) handleClient(socket *websocket.Conn, variables map[string]string) {
	inData := make(chan network.Packet)
	outData := make(chan network.Packet)
	closed := make(chan struct{})
	connection := network.Connection{
		InData:  inData,
		OutData: outData,
		Closed:  closed,
	}
	wsConn := WsConnection{socket, connection}
	endpoint.connections[wsConn] = struct{}{}
	endpoint.OnConnect <- network.NewConnection{
		Connection: connection,
		Variables:  variables,
	}

	go func() {
		for pkt := range outData {
			socket.WriteMessage(websocket.BinaryMessage, pkt.Data)
		}
		socket.Close()
	}()

	err := <-socketReceiveLoop(wsConn, inData)
	close(closed)
	fmt.Println("Client disconnected ", err)
	delete(endpoint.connections, wsConn)
}

func socketReceiveLoop(connection WsConnection, onData chan<- network.Packet) chan error {
	onError := make(chan error)
	go func() {
		for {
			_, data, err := connection.socket.ReadMessage()
			if err != nil {
				onError <- err
				break
			}
			fmt.Println("Data received WsEndpoint.go")
			var message = network.Packet{data, connection.connection}
			onData <- message
		}
	}()
	return onError
}

func StartWebSocketServer(port int, paths []string) (WsEndpoint, error) {
	endpoint := WsEndpoint{
		make(chan network.NewConnection),
		make(map[WsConnection]struct{}),
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(_ *http.Request) bool { return true },
	}

	handlerFunc := func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		socket, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print(err)
			return
		}
		endpoint.handleClient(socket, vars)
	}

	router := mux.NewRouter()
	for _, path := range paths {
		router.HandleFunc(path, handlerFunc)
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
			log.Fatal("Unable to start listen on http server: ", err)
		}
	}()
	return endpoint, nil
}
