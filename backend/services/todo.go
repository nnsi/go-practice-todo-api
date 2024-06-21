package services

import (
	"encoding/json"
	"errors"
	"go-practice-todo/infra"
	"go-practice-todo/models"
	"go-practice-todo/repositories"
	"log"
)

type TodoService struct {
	repo     repositories.TodoRepositoryInterface
	Notifier infra.NotifierInterface
}

func NewTodoService(repo repositories.TodoRepositoryInterface, notifier infra.NotifierInterface) *TodoService {
	return &TodoService{repo: repo, Notifier: notifier}
}

func (s *TodoService) Index(isShowDeleted bool, userID string) ([]models.Todo, error) {
	if userID == "" {
		return nil, errors.New("user not found")
	}
	return s.repo.FindAll(isShowDeleted, userID)
}

func (s *TodoService) Show(id string, userID string) (*models.Todo, error) {
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

	data, _err := json.Marshal(todo)
	if _err != nil {
		return todo, _err
	}

	log.Println("Created: ", string(data))
	s.Notifier.BroadcastMessage(infra.Message{
		Event:  "create",
		Data:   string(data),
		UserID: userID,
	})
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

	data, _err := json.Marshal(todo)
	if _err != nil {
		return todo, _err
	}

	log.Println("Updated:: ", string(data))
	s.Notifier.BroadcastMessage(infra.Message{
		Event:  "update",
		Data:   string(data),
		UserID: userID,
	})

	return todo, err
}

func (s *TodoService) Delete(id string, userID string) error {
	err := s.repo.Delete(id, userID)

	data, _err := json.Marshal(models.Todo{ID: id})
	if _err != nil {
		return err
	}

	log.Println("Deleted: ", string(data))
	s.Notifier.BroadcastMessage(infra.Message{
		Event:  "delete",
		Data:   string(data),
		UserID: userID,
	})

	return err
}
