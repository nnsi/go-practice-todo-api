package repositories

import (
	"go-practice-todo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TodoRDBRepository struct {
	db *gorm.DB
}

func NewTodoRDBRepository(dsn string) (*TodoRDBRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// マイグレーションを実行してテーブルを作成
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		return nil, err
	}

	return &TodoRDBRepository{db: db}, nil
}

func (r *TodoRDBRepository) Index() ([]models.Todo, error) {
	var todos []models.Todo
	result := r.db.Find(&todos)
	return todos, result.Error
}

func (r *TodoRDBRepository) Show(id string) (*models.Todo, error) {
	var todo models.Todo
	result := r.db.First(&todo, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func (r *TodoRDBRepository) Create(todo *models.Todo) error {
	result := r.db.Create(todo)
	return result.Error
}

func (r *TodoRDBRepository) Update(todo *models.Todo) error {
	result := r.db.Save(todo)
	return result.Error
}

func (r *TodoRDBRepository) Delete(id string) error {
	result := r.db.Delete(&models.Todo{}, "id = ?", id)
	return result.Error
}
