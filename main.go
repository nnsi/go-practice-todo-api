package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

func generateULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("GET /todos", Index)
	http.HandleFunc("OPTIONS /todos", requestOptions)
	http.HandleFunc("POST /todos", Create)
	http.HandleFunc("GET /todos/{id}", Show)
	http.HandleFunc("PUT /todos/{id}", Update)
	http.HandleFunc("DELETE /todos/{id}", Delete)

	http.ListenAndServe("localhost:8080", nil)
}

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo = []Todo{
	{ID: "01HZW7K57BWVW46K95AAHV3H7Q", Title: "Buy milk", Completed: false},
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	json.NewEncoder(w).Encode(data)
}

func requestOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

func Index(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, todos)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	var fields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		writeJSONResponse(w, "Invalid JSON")
		return
	}
	todo.ID = generateULID()

	if title, ok := fields["title"].(string); ok {
		todo.Title = title
	} else {
		writeJSONResponse(w, "Title is required")
		return
	}
	todo.Completed = false
	todos = append(todos, todo)
	writeJSONResponse(w, todo)
}

func Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for i, todo := range todos {
		if todo.ID == id {
			writeJSONResponse(w, todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedFields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			if title, ok := updatedFields["title"].(string); ok {
				todos[i].Title = title
			}
			if completed, ok := updatedFields["completed"].(bool); ok {
				todos[i].Completed = completed
			}
			writeJSONResponse(w, todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			writeJSONResponse(w, todos)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}
