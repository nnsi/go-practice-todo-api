package repositories

import (
	"errors"
	"go-practice-todo/models"
)

type TodoRepository struct {
	todos map[string]models.Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos: make(map[string]models.Todo),
	}
}

func (r *TodoRepository) Index() ([]models.Todo, error) {
	todos := []models.Todo{}
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) Show(id string) (*models.Todo, error) {
	for i, todo := range r.todos {
		if todo.ID == id {
			temp := r.todos[i]
			return &temp, nil
		}
	}
	return nil, errors.New("todo not found")
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	r.todos[todo.ID] = *todo
	return nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	for i, t := range r.todos {
		if t.ID == todo.ID {
			r.todos[i] = *todo
			return nil
		}
	}
	return errors.New("todo not found")
}

func (r *TodoRepository) Delete(id string) error {
	for _, todo := range r.todos {
		if todo.ID == id {
			delete(r.todos, id)
			return nil
		}
	}
	return errors.New("todo not found")
}
