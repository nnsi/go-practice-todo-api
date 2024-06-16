package services

import (
	"time"

	"go-practice-todo/models"
	"go-practice-todo/repositories"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repositories.UserRepositoryInterface
	jwtSecret string
}

func NewUserService(repo repositories.UserRepositoryInterface, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Create(userDto *models.UserDTO) (*models.User, error) {
	user := &models.User{
		LoginID: userDto.LoginID,
		Name:    userDto.Username,
		Password: userDto.Password,
	}
	user.ID = GenerateULID()
	err := s.repo.Create(user)
	return user, err
}

func (s *UserService) FindByID(loginID string) (*models.User, error) {
	return s.repo.FindByID(loginID)
}

func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Name,
		"exp":      time.Now().Add(time.Hour * 24 * 365).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	return token, err
}