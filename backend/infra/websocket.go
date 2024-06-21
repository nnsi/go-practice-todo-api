package infra

import (
	"log"
	"sync"
	"time"

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
		broadcast: make(chan Message, 100),
	}
}

func (n *WebSocketNotifier) RegisterClient(userID string, ws *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.clients[userID] == nil {
		n.clients[userID] = make(map[*websocket.Conn]bool)
	}
	n.clients[userID][ws] = true
	log.Printf("Registered client: %s, User ID: %s", ws.RemoteAddr(), userID)
}

func (n *WebSocketNotifier) UnregisterClient(userID string, ws *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.clients[userID][ws]; ok {
		delete(n.clients[userID], ws)
		log.Printf("Unregistered client: %s, User ID: %s", ws.RemoteAddr(), userID)
		if len(n.clients[userID]) == 0 {
			delete(n.clients, userID)
		}
	} else {
		log.Printf("Attempted to unregister a non-registered client: %s, User ID: %s", ws.RemoteAddr(), userID)
	}
}

func (n *WebSocketNotifier) BroadcastMessage(msg Message) {
	go func() {
		n.broadcast <- msg
	}()
}

func (n *WebSocketNotifier) Start() {
	go func() {
		for {
			msg := <-n.broadcast
			log.Printf("Broadcasting message: %v", msg)
			//			n.mu.Lock()
			for client := range n.clients[msg.UserID] {
				go func(client *websocket.Conn) {
					log.Printf("Sending message to client: %s", client.RemoteAddr())

					if err := client.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
						log.Printf("Error setting write deadline: %v", err)
					}

					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("Error writing JSON to client: %v", err)
						client.Close()
						delete(n.clients[msg.UserID], client)
						if len(n.clients[msg.UserID]) == 0 {
							delete(n.clients, msg.UserID)
						}
					} else {
						log.Printf("Message successfully sent to client: %s", client.RemoteAddr())
					}
				}(client)
			}
			//			n.mu.Unlock()
		}
	}()
}
