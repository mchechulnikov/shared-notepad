package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: 	 func(r *http.Request) bool { return true },
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	currentConnection, err := upgrader.Upgrade(w, r, nil)
	if err  != nil {
		return
	}

	initialValue := struct {
		roomEvent
		textChangedEvent
	}{}
	initialValue.Text = room.Text
	initialValueJSON, _ := json.Marshal(initialValue)
	if err = currentConnection.WriteMessage(1, initialValueJSON); err != nil {
		currentConnection.Close()
		return
	}

	clientsToRoomsMapping[Client{currentConnection}] = room

	for {
		msgType, msg, err := currentConnection.ReadMessage()
		if err != nil {
			return
		}

		if err := handleMessage(msg); err != nil {
			log.Println(err)
			continue
		}

		// broadcast
		for client := range clientsToRoomsMapping {
			if client.Connection == currentConnection {
				continue
			}

			if err = client.Connection.WriteMessage(msgType, msg); err != nil {
				client.Connection.Close()
				delete(clientsToRoomsMapping, client)
			}
		}
	}
}

func handleMessage(msg []byte) error {
	var roomEvent roomEvent
	buffer := make([]byte, len(msg))
	copy(buffer, msg)
	if err := json.Unmarshal(buffer, &roomEvent); err != nil {
		return fmt.Errorf("unmarshalling WS message error: %s", err)
	}

	var unmarshalErr error
	switch roomEvent.EventType {
	case textChangedEventType:
		var textChangedEvent textChangedEvent
		unmarshalErr = json.Unmarshal(buffer, &textChangedEvent)
		saveText(textChangedEvent.Text)
		break
	case cursorPositionChangedEventType:
		var cursorPositionEvent cursorPositionEvent
		unmarshalErr = json.Unmarshal(buffer, &cursorPositionEvent)
		saveCursorPosition(cursorPositionEvent.CursorPosition, roomEvent.ActorName)
		break
	case cursorPositionCancellationEventType:
		saveCursorCancellation(roomEvent.ActorName)
		break
	case selectionChangedEventType:
		var selectionEvent selectionEvent
		unmarshalErr = json.Unmarshal(buffer, &selectionEvent)
		saveSelection(selectionEvent.SelectionStart, selectionEvent.SelectionEnd, roomEvent.ActorName)
		break
	case selectionCancelledEventType:
		saveSelectionCancellation(roomEvent.ActorName)
		break
	}
	if unmarshalErr != nil {
		return fmt.Errorf("unmarshalling data received from message error: %s", unmarshalErr)
	}

	return nil
}

func saveText(text string) {
	room.Text = text
}

func saveCursorPosition(cursorPosition int, actor string) {
	room.CursorsPositions[actor] = CursorPosition{ActorName: actor, Position: cursorPosition}
}

func saveSelection(start int, end int, actor string) {
	room.Selections[actor] = Selection{ActorName: actor, Start: start, End: end}
}

func saveCursorCancellation(actor string) {
	delete(room.CursorsPositions, actor)
}

func saveSelectionCancellation(actor string) {
	delete(room.Selections, actor)
}

const (
	textChangedEventType = "TextChanged"
	cursorPositionChangedEventType = "CursorPositionChanged"
	cursorPositionCancellationEventType = "CursorPositionCancelled"
	selectionChangedEventType = "SelectionChanged"
	selectionCancelledEventType = "SelectionCancelled"
)

type roomEvent struct {
	ActorName string `json:"actor_name"`
	EventType string `json:"event_type"`
}

type textChangedEvent struct {
	Text string `json:"text"`
}

type cursorPositionEvent struct {
	CursorPosition int `json:"cursor_position"`
}

type selectionEvent struct {
	SelectionStart int `json:"selection_start"`
	SelectionEnd int `json:"selection_end"`
}