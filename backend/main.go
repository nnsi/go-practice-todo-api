package main

import (
	"go-practice-todo/handlers"
	"go-practice-todo/infra"
	"go-practice-todo/repositories"
	"go-practice-todo/services"
	"log"
	"net/http"
)

const JWT_SECRET = "secret"

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=todoapp port=5432 TimeZone=Asia/Tokyo"

	notifier := infra.NewWebSocketNotifier()
	notifier.Start()

	todoRepo, err := repositories.NewTodoRDBRepository(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	todoService := services.NewTodoService(todoRepo, notifier)
	todoHandler := handlers.NewTodoHandler(todoService)

	userRepo, err := repositories.NewUserRepository(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	userService := services.NewUserService(userRepo, JWT_SECRET)
	authHandler := handlers.NewAuthHandler(userService)

	wsHandler := handlers.NewWebSocketHandler(notifier, todoService)

	Routes(todoHandler, authHandler, wsHandler, userService)

	http.ListenAndServe("localhost:8080", nil)
}
