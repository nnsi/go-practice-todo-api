package repositories

import "go-practice-todo/models"

// TodoRepository defines the methods that any
// data storage provider needs to implement to get
// and store todos.
type UserRepositoryInterface interface {
	FindByID(id string) (*models.User, error)
	Create(user *models.User) error
}
