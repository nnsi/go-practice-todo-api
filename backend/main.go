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
	todoRepo, err := repositories.NewTodoRDBRepository (dsn) 
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	userRepo, err := repositories.NewUserRepository(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	userService := services.NewUserService(userRepo, "secret")
	authHandler := handlers.NewAuthHandler(userService)
	
	Routes(todoHandler, authHandler)

	http.ListenAndServe("localhost:8080", nil)
}
