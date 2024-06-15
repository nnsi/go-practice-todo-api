package repositories

import "go-practice-todo/models"

// TodoRepository defines the methods that any
// data storage provider needs to implement to get
// and store todos.
type TodoRepositoryInterface interface {
	FindAll(isShowDeleted bool) ([]models.Todo, error)
	FindByID(id string) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id string) error
}