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
	OnData       chan network.Packet
	OnConnect    chan network.NewConnection
	OnDisconnect chan network.Connection

	connections map[WsConnection]bool
}

type WsConnection struct {
	socket *websocket.Conn
}

func (con WsConnection) Send(data []byte) error {
	return con.socket.WriteMessage(websocket.BinaryMessage, data)
}

func (con WsConnection) Close() {
	con.socket.Close()
}

func (endpoint *WsEndpoint) clientDisconnect(connection WsConnection) {
	delete(endpoint.connections, connection)
	endpoint.OnDisconnect <- connection
}

func (endpoint *WsEndpoint) handleClient(connection *websocket.Conn, variables map[string]string) {
	wsConn := WsConnection{connection}
	endpoint.connections[wsConn] = true
	endpoint.OnConnect <- network.NewConnection{
		Connection: wsConn,
		Variables:  variables,
	}

	err := <-socketReceiveLoop(wsConn, endpoint.OnData)
	fmt.Println("Client disconnected ", err)
	endpoint.clientDisconnect(wsConn)
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
			var message = network.Packet{data, connection}
			onData <- message
		}
	}()
	return onError
}

func StartWebSocketServer(port int, paths []string) (WsEndpoint, error) {
	endpoint := WsEndpoint{
		make(chan network.Packet),
		make(chan network.NewConnection),
		make(chan network.Connection),
		make(map[WsConnection]bool),
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
