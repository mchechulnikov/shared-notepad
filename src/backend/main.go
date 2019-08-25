package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: 	 func(r *http.Request) bool { return true },
}

var room = &Room {
	uuid.New(),
	`import fmt

func main() {
	fmt.Printf("Hello world!")
}`,
}

var clientsToRoomsMapping = make(map[Client]*Room, 0)

func main() {
	fmt.Println("App started")
	http.HandleFunc("/room/join", WebsocketHandler)
	_ = http.ListenAndServe(":5000", nil)
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	currentConnection, err := upgrader.Upgrade(w, r, nil)
	if err  != nil {
		return
	}

	if err = currentConnection.WriteMessage(1, []byte(room.text)); err != nil {
		currentConnection.Close()
		return
	}

	clientsToRoomsMapping[Client{currentConnection}] = room

	for {
		msgType, msg, err := currentConnection.ReadMessage()
		if err != nil {
			return
		}

		room.text = string(msg)

		for client := range clientsToRoomsMapping {
			if client.connection == currentConnection {
				continue
			}

			if err = client.connection.WriteMessage(msgType, msg); err != nil {
				client.connection.Close()
				delete(clientsToRoomsMapping, client)
			}
		}
	}
}

type Client struct {
	connection	*websocket.Conn
}

type Room struct {
	id		uuid.UUID
	text 	string
}