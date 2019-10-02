package server

import (
	"encoding/base64"
	"github.com/google/uuid"
)

var room = &Room {
	uuid.New(),
	base64.StdEncoding.EncodeToString([]byte(`import fmt

func main() {
	fmt.Printf("Hello world!")
}`)),
	make(map[string]CursorPosition, 0),
	make(map[string]Selection, 0),
}

var clientsToRoomsMapping = make(map[Client]*Room, 0)