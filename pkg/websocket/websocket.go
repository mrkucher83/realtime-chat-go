package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrader require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ws, nil
}
