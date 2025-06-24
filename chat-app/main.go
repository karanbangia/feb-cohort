package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub manages active WebSocket connections
type Hub struct {
	connectionMap map[*websocket.Conn]bool
}

// NewHub initializes a new Hub
func NewHub() *Hub {
	return &Hub{
		connectionMap: make(map[*websocket.Conn]bool),
	}
}

// Register a new connection
func (h *Hub) register(conn *websocket.Conn) {
	h.connectionMap[conn] = true
}

// Unregister a connection and close it
func (h *Hub) unregister(conn *websocket.Conn) {
	if _, exists := h.connectionMap[conn]; exists {
		delete(h.connectionMap, conn)
		conn.Close()
	}
}

// Broadcast message to all clients except sender
func (h *Hub) broadcast(sender *websocket.Conn, redisClient *redis.Client, msg []byte) {

	for conn := range h.connectionMap {
		if conn != sender {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Error writing message:", err)
				h.unregister(conn) // Remove faulty connection
			}
			redisClient.Publish(context.Background(), "chat", msg)
		}
	}
}

// Handle new WebSocket connections
func handleConnections(h *Hub, redisClient *redis.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	h.register(conn)

	// Handle incoming messages
	go func() {
		defer h.unregister(conn) // Ensure cleanup on exit

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Connection closed by client")
				} else {
					log.Println("Read error:", err)
				}
				break
			}
			h.broadcast(conn, redisClient, msg)
		}
	}()
}

func main() {
	hub := NewHub()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, redisClient, w, r)
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
