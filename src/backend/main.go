package main

import (
	"fmt"
	"net/http"
	"sharednotepad/server"
)

func main() {
	fmt.Println("App started")
	http.HandleFunc("/room/join", server.WSHandler)
	_ = http.ListenAndServe(":5000", nil)
}