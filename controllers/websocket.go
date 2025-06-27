package controllers

import (
	"log"

	"website-builder/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(c *gin.Context, hub *websocket.Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	// Get user from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		conn.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
		conn.Close()
		return
	}

	client := &websocket.Client{
		Conn:   conn,
		UserID: userID.(string),
		Send:   make(chan []byte, 256),
	}

	hub.Register <- client

	// Start goroutines for reading and writing messages
	go client.WritePump()
	go client.ReadPump(hub)
}
