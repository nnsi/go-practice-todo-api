package services

import (
	"math/rand"
	"time"

	"go-practice-todo/models"
	"go-practice-todo/repositories"

	"github.com/oklog/ulid/v2"
)

type TodoService struct {
	repo repositories.TodoRepositoryInterface
}

func NewTodoService(repo repositories.TodoRepositoryInterface) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) generateULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func (s *TodoService) Index() ([]models.Todo, error) {
	return s.repo.Index()
}

func (s *TodoService) Show(id string) (*models.Todo, error) {
	return s.repo.Show(id)
}

func (s *TodoService) Create(title string) (*models.Todo, error) {
	todo := &models.Todo{
		ID:        s.generateULID(),
		Title:     title,
		Completed: false,
	}
	err := s.repo.Create(todo)
	return todo, err
}

func (s *TodoService) Update(id string, title string, completed bool) (*models.Todo, error) {
	todo, err := s.repo.Show(id)
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