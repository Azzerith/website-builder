package websocket

import "log"

type Client struct {
    send chan []byte
    // Add other client fields as needed
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Client connected. Total clients: %d", len(h.clients))
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Client disconnected. Total clients: %d", len(h.clients))
            }
        case message := <-h.broadcast:
            log.Printf("Broadcasting message to %d clients", len(h.clients))
            for client := range h.clients {
                select {
                case client.send <- message:
                    // Message sent successfully
                default:
                    close(client.send)
                    delete(h.clients, client)
                    log.Println("Client disconnected due to slow connection")
                }
            }
        }
    }
}