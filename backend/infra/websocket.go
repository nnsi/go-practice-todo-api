package infra

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketNotifier struct {
	clients   map[string]map[*websocket.Conn]bool
	broadcast chan Message
	mu        sync.Mutex
}

type Message struct {
	Event  string `json:"event"`
	Data   string `json:"data"`
	UserID string `json:"user_id"`
}

func NewWebSocketNotifier() *WebSocketNotifier {
	return &WebSocketNotifier{
			clients:   make(map[string]map[*websocket.Conn]bool),
			broadcast: make(chan Message),
	}
}

func (n *WebSocketNotifier) RegisterClient(userID string, ws *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.clients[userID] == nil {
			n.clients[userID] = make(map[*websocket.Conn]bool)
	}
	n.clients[userID][ws] = true
}

func (n *WebSocketNotifier) UnregisterClient(userID string, ws *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.clients[userID], ws)
	if len(n.clients[userID]) == 0 {
			delete(n.clients, userID)
	}
}

func (n *WebSocketNotifier) BroadcastMessage(msg Message) {
	n.broadcast <- msg
}

func (n *WebSocketNotifier) Start() {
	go func() {
			for {
					msg := <-n.broadcast
					n.mu.Lock()
					for client := range n.clients[msg.UserID] {
							err := client.WriteJSON(msg)
							if err != nil {
									client.Close()
									delete(n.clients[msg.UserID], client)
									if len(n.clients[msg.UserID]) == 0 {
											delete(n.clients, msg.UserID)
									}
							}
					}
					n.mu.Unlock()
			}
	}()
}