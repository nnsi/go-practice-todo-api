package handlers

import (
	"encoding/json"
	"go-practice-todo/infra"
	"go-practice-todo/services"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Notifier    infra.NotifierInterface
	TodoService *services.TodoService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketHandler(notifier infra.NotifierInterface, todoService *services.TodoService) *WebSocketHandler {
	return &WebSocketHandler{
		Notifier:    notifier,
		TodoService: todoService,
	}
}

func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims["user_id"].(string)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	h.Notifier.RegisterClient(userID, ws)

	for {
		var msg infra.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			h.Notifier.UnregisterClient(userID, ws)
			break
		}
		// メッセージのイベントに基づいて処理
		switch msg.Event {
		case "get_todos":
			todos, err := h.TodoService.Index(false, userID)
			if err != nil {
				log.Printf("error getting todos: %v", err)
				break
			}

			data, _err := json.Marshal(todos)
			if _err != nil {
				return
			}

			response := infra.Message{
				Event:  "list",
				Data:   string(data),
				UserID: userID,
			}
			err = ws.WriteJSON(response)
			if err != nil {
				log.Printf("error sending todos: %v", err)
			}
		default:
			log.Printf("unknown event: %v", msg.Event)
		}
	}
}
