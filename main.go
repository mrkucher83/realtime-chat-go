package main

import (
	"fmt"
	"github.com/mrkucher83/realtime-chat-go/pkg/websocket"
	"log"
	"net/http"
)

// define our WebSocket endpoint
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: ws,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01: http://localhost:8080")
	setupRoutes()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start the web server - Error: %v", err)
	}
}
