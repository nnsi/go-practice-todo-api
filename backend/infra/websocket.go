package infra

import (
	"log"
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
	log.Printf("Registered client: %s, User ID: %s, Clients: %d", ws.RemoteAddr(), userID, len(n.clients[userID]))
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
		for msg := range n.broadcast {
			msg := msg
			n.mu.Lock()

			clients, ok := n.clients[msg.UserID]
			if !ok {
				log.Printf("No clients found for User ID: %s", msg.UserID)
				n.mu.Unlock()
				continue
			}

			// ローカルコピーを作成してロックを早期に解放
			clientList := make([]*websocket.Conn, 0, len(clients))
			for client := range clients {
				clientList = append(clientList, client)
			}
			n.mu.Unlock()
			log.Printf("Broadcasting message to %d clients", len(clientList))

			for _, client := range clientList {
				go func(client *websocket.Conn) {
					log.Printf("Sending message to client: %s", client.RemoteAddr())

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
		}
	}()
}
