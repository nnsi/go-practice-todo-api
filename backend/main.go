package main

import (
	"fmt"
	"go-practice-todo/handlers"
	"go-practice-todo/infra"
	"go-practice-todo/repositories"
	"go-practice-todo/services"
	"log"
	"net/http"
	"os"
)

const JWT_SECRET = "secret"

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {

	DB_HOST := getEnv("DB_HOST", "db")
	DB_USER := getEnv("DB_USER", "postgres")
	DB_PASSWORD := getEnv("DB_PASSWORD", "postgres")
	DB_NAME := getEnv("DB_NAME", "todoapp")
	DB_PORT := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Tokyo", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

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

	http.ListenAndServe("0.0.0.0:8080", nil)
}
