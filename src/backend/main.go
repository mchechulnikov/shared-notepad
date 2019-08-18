package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: 	 func(r *http.Request) bool { return true },
}

var connections = make([]*websocket.Conn, 0)

func main() {
	fmt.Println("Shared notepad app started")
	http.HandleFunc("/room/join", WebsocketHandler)

	_ = http.ListenAndServe(":5000", nil)
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	currentConnection, err := upgrader.Upgrade(w, r, nil)
	if err  != nil {
		return
	}
 	connections = append(connections, currentConnection)

	for {
		msgType, msg, err := currentConnection.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("%s received: %s\n", currentConnection.RemoteAddr(), string(msg))

		for _, connection := range connections {
			if connection == currentConnection {
				continue
			}

			if err = connection.WriteMessage(msgType, msg); err != nil {
				return
			}

			fmt.Printf("%s sent: %s\n", connection.RemoteAddr(), string(msg))
		}
	}
}