package repositories

import "go-practice-todo/models"

// TodoRepository defines the methods that any
// data storage provider needs to implement to get
// and store todos.
type TodoRepositoryInterface interface {
	FindAll(isShowDeleted bool, userID string) ([]models.Todo, error)
	FindByID(id string, userID string) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id string, userID string) error
}
