package main

import (
	"go-practice-todo/handlers"
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

func Routes(todoHandler *handlers.TodoHandler, authHandler *handlers.AuthHandler) {

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	http.HandleFunc("OPTIONS /todos", RequestOptions)

	http.Handle("GET /todos", ChainMiddleware(http.HandlerFunc(todoHandler.Index), authHandler.AuthMiddleware))
	http.Handle("POST /todos", ChainMiddleware(http.HandlerFunc(todoHandler.Create), authHandler.AuthMiddleware))
	http.Handle("GET /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Show), authHandler.AuthMiddleware))
	http.Handle("PUT /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Update), authHandler.AuthMiddleware))
	http.Handle("DELETE /todos/{id}", ChainMiddleware(http.HandlerFunc(todoHandler.Delete), authHandler.AuthMiddleware))
}
