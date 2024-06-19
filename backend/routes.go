package main

import (
	"go-practice-todo/handlers"
	"go-practice-todo/middleware"
	"go-practice-todo/services"

	"net/http"
)

func RequestOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func Routes(todoHandler *handlers.TodoHandler, authHandler *handlers.AuthHandler, wsHandler *handlers.WebSocketHandler, userService *services.UserService) {

	http.HandleFunc("OPTIONS /ws", RequestOptions)
	http.HandleFunc("GET /ws", wsHandler.HandleConnections)

	http.HandleFunc("OPTIONS /register", RequestOptions)
	http.HandleFunc("POST /register", authHandler.Register)

	http.HandleFunc("OPTIONS /login", RequestOptions)
	http.HandleFunc("POST /login", authHandler.Login)

	http.HandleFunc("OPTIONS /todos", RequestOptions)
	http.HandleFunc("OPTIONS /todos/", RequestOptions)

	http.Handle("GET /todos", ChainMiddleware(http.HandlerFunc(todoHandler.Index), middleware.AuthMiddleware(userService)))
	http.Handle("POST /todos", ChainMiddleware(http.HandlerFunc(todoHandler.Create), middleware.AuthMiddleware(userService)))
	http.Handle("GET /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Show), middleware.AuthMiddleware(userService)))
	http.Handle("PUT /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Update), middleware.AuthMiddleware(userService)))
	http.Handle("DELETE /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Delete), middleware.AuthMiddleware(userService)))
}
