package handlers

import (
	"encoding/json"
	"go-practice-todo/services"
	"net/http"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	isShowDeleted := r.URL.Query().Get("show-deleted") == "1"
	todos, err := h.service.Index(isShowDeleted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSONResponse(w, todos)
}

func (h *TodoHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := h.service.Show(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	WriteJSONResponse(w, todo)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var fields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		WriteJSONResponse(w, "Invalid JSON")
		return
	}

	title, ok := fields["title"].(string)
	if !ok {
		WriteJSONResponse(w, "Title is required")
		return
	}

	todo, err := h.service.Create(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSONResponse(w, todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		todo, err := h.service.Show(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var fields map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			WriteJSONResponse(w, "Invalid JSON")
			return
		}

		if title, ok := fields["title"].(string); ok {
			todo.Title = title
		}
		if completed, ok := fields["completed"].(bool); ok {
			todo.Completed = completed
		}
		todo, err = h.service.Update(id, todo.Title, todo.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		WriteJSONResponse(w, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	WriteJSONResponse(w, map[string]string{"id": id})
}