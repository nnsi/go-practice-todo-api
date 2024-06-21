package infra

import "github.com/gorilla/websocket"

type Message struct {
	Event  string `json:"event"`
	Data   string `json:"data"`
	UserID string `json:"user_id"`
}

type NotifierInterface interface {
	BroadcastMessage(msg Message)
	RegisterClient(userID string, ws *websocket.Conn)
	Start()
	UnregisterClient(userID string, ws *websocket.Conn)
}
