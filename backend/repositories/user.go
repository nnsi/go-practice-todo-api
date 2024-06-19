package repositories

import (
	"go-practice-todo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dsn string) (*UserRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// マイグレーションを実行してテーブルを作成
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) FindByID(loginId string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("login_id = ?", loginId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}
