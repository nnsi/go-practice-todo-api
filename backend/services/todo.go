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

func (s *TodoService) Index(isShowDeleted bool) ([]models.Todo, error) {
	return s.repo.FindAll(isShowDeleted)
}

func (s *TodoService) Show(id string) (*models.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *TodoService) Create(title string) (*models.Todo, error) {
	todo := &models.Todo{
		ID:        GenerateULID(),
		Title:     title,
		Completed: false,
	}
	err := s.repo.Create(todo)
	return todo, err
}

func (s *TodoService) Update(id string, title string, completed bool) (*models.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	todo.Title = title
	todo.Completed = completed
	err = s.repo.Update(todo)
	return todo, err
}

func (s *TodoService) Delete(id string) error {
	return s.repo.Delete(id)
}