package infra

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type RedisNotifier struct {
	Client    *redis.Client
	Channel   string
	Broadcast chan Message
	ctx       context.Context
	clients   map[string]map[*websocket.Conn]bool
	mu        sync.Mutex
}

func NewRedisNotifier(redisURL string, channel string) *RedisNotifier {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	return &RedisNotifier{
		Broadcast: make(chan Message, 100),
		Client:    client,
		Channel:   channel,
		clients:   make(map[string]map[*websocket.Conn]bool),
		ctx:       context.Background(),
	}
}

func (n *RedisNotifier) Start() {
	go n.subscribeToChannel()
}

func (n *RedisNotifier) subscribeToChannel() {
	pubsub := n.Client.Subscribe(n.ctx, n.Channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(n.ctx)
		if err != nil {
			log.Printf("failed to receive message: %v", err)
			continue
		}

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		n.broadcastToClients(message)
	}
}

func (n *RedisNotifier) broadcastToClients(msg Message) {
	n.mu.Lock()
	clients, ok := n.clients[msg.UserID]
	if !ok {
		log.Printf("No clients found for User ID: %s", msg.UserID)
		n.mu.Unlock()
		return
	}

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
				n.mu.Lock()
				delete(n.clients[msg.UserID], client)
				if len(n.clients[msg.UserID]) == 0 {
					delete(n.clients, msg.UserID)
				}
				n.mu.Unlock()
			} else {
				log.Printf("Message successfully sent to client: %s", client.RemoteAddr())
			}
		}(client)
	}
}

func (n *RedisNotifier) BroadcastMessage(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}
	err = n.Client.Publish(n.ctx, n.Channel, data).Err()
	if err != nil {
		log.Printf("failed to publish message: %v", err)
	}
}

func (n *RedisNotifier) RegisterClient(userID string, ws *websocket.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.clients[userID] == nil {
		n.clients[userID] = make(map[*websocket.Conn]bool)
	}
	n.clients[userID][ws] = true
	log.Printf("Registered client: %s, User ID: %s, Clients: %d", ws.RemoteAddr(), userID, len(n.clients[userID]))
}

func (n *RedisNotifier) UnregisterClient(userID string, ws *websocket.Conn) {
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
