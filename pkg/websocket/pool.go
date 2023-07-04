package websocket

import (
	"fmt"
	"log"
)

// Pool is a struct which will contain all the channels we need for concurrent communication, as well as a map of clients.
type Pool struct {
	Register   chan *Client     // notify all the clients within this pool when a new client connects.
	Unregister chan *Client     // unregister a user and notify the pool when a client disconnects.
	Clients    map[*Client]bool // a map uses the boolean value to dictate active/inactive.
	Broadcast  chan Message     // loop through all clients in the pool and send the message through the socket connection.
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection Pool: ", len(pool.Clients))
			notifyUsers("New User Joined...", pool)
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of connection Pool: ", len(pool.Clients))
			notifyUsers("User Disconnected...", pool)
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending a message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Printf("Failed to notify users: %v", err)
				}
			}

		}
	}
}

func notifyUsers(text string, pool *Pool) {
	for client := range pool.Clients {
		if err := client.Conn.WriteJSON(Message{1, text}); err != nil {
			log.Printf("Failed to notify users: %v", err)
		}
	}
}
