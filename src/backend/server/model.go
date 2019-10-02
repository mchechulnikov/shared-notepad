package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Connection	*websocket.Conn
}

type Room struct {
	Id		uuid.UUID
	Text 	string
	CursorsPositions map[string]CursorPosition
	Selections map[string]Selection
}

type CursorPosition struct {
	ActorName string
	Position int
}

type Selection struct {
	ActorName string
	Start int
	End int
}