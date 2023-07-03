package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrader require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// we'll need to check the origin of our connection
	// this will allow to make requests from React development server to here.
	// For now, it just allows any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// reader will listen for new messages being sent to our WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Listen indefinitely for new messages coming through on our WevSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Simple Server"); err != nil {
			log.Printf("Failed to handle request - Error: %v", err)
		}
	})

	http.HandleFunc("/ws", serveWs)
}

func main() {
	fmt.Println("Chat App v0.01: http://localhost:8080")
	setupRoutes()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start the web server - Error: %v", err)
	}
}
