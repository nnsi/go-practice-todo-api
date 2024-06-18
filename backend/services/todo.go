package services

import (
	"go-practice-todo/models"
	"go-practice-todo/repositories"
)

type TodoService struct {
	repo repositories.TodoRepositoryInterface
}

func NewTodoService(repo repositories.TodoRepositoryInterface) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) Index(isShowDeleted bool,userID string) ([]models.Todo, error) {
	return s.repo.FindAll(isShowDeleted, userID)
}

func (s *TodoService) Show(id string,userID string) (*models.Todo, error) {
	return s.repo.FindByID(id, userID)
}

func (s *TodoService) Create(title string, userID string) (*models.Todo, error) {
	todo := &models.Todo{
		ID:        GenerateULID(),
		Title:     title,
		Completed: false,
		UserID:    userID,
	}
	err := s.repo.Create(todo)
	return todo, err
}

func (s *TodoService) Update(id string, title string, completed bool, userID string) (*models.Todo, error) {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	todo.Title = title
	todo.Completed = completed
	err = s.repo.Update(todo)
	return todo, err
}

func (s *TodoService) Delete(id string,userID string) error {
	return s.repo.Delete(id, userID)
}