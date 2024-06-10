package main

import (
	"go-practice-todo/handlers"
	"go-practice-todo/repositories"
	"go-practice-todo/services"
	"log"
	"net/http"
)

func main() { 
	// repo := repositories.NewTodoRepository()
	dsn := "host=localhost user=postgres password=postgres dbname=todoapp port=5432 TimeZone=Asia/Tokyo"
	repo, err := repositories.NewTodoRDBRepository (dsn) 
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)
	Routes(handler)

	http.ListenAndServe("localhost:8080", nil)
}
