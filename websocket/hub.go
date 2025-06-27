package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	Send   chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client disconnected. Total clients: %d", len(h.Clients))
			}
		case message := <-h.Broadcast:
			log.Printf("Broadcasting message to %d clients", len(h.Clients))
			for client := range h.Clients {
				select {
				case client.Send <- message:
					// Message sent successfully
				default:
					close(client.Send)
					delete(h.Clients, client)
					log.Println("Client disconnected due to slow connection")
				}
			}
		}
	}
}

func (c *Client) ReadPump(hub *Hub) {
	// isi sesuai kebutuhan, contoh sederhana:
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			hub.Unregister <- c
			c.Conn.Close()
			break
		}
		hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.Conn.Close()
}
