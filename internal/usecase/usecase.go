package usecase

import (
	"errors"
	"strings"

	"github.com/SaveliiYam/http_service_with_db/internal/models"
	"github.com/SaveliiYam/http_service_with_db/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (models.User, error) {
	user, ok := s.repo.GetByID(id)
	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) CreateUser(name string) (models.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.User{}, ErrInvalidUser
	}

	return s.repo.Create(name), nil
}
