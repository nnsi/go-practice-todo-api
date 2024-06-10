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

func Routes(handler *handlers.TodoHandler) {
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("GET /todos", handler.Index)
	http.HandleFunc("OPTIONS /todos", RequestOptions)
	http.HandleFunc("POST /todos", handler.Create)
	http.HandleFunc("GET /todos/{id}", handler.Show)
	http.HandleFunc("PUT /todos/{id}", handler.Update)
	http.HandleFunc("DELETE /todos/{id}", handler.Delete)
}