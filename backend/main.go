package main

import (
	"go-practice-todo/handlers"
	"go-practice-todo/repositories"
	"go-practice-todo/services"
	"net/http"
)

func main() {
	repo := repositories.NewTodoRepository()
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)
	Routes(handler)

	http.ListenAndServe("localhost:8080", nil)
}
