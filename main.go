package main

import (
	"fmt"
	"net/http"
)

func main() {
	//http.HandleFunc("/", homePage)
	//http.HandleFunc("/ws", handleConnections)

	//go handleMessages()
	manager := NewManager()
	http.HandleFunc("/ws", manager.wsHandler)
	fmt.Println("Websocket Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
