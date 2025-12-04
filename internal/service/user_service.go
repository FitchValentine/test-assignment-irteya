package service

import (
	"errors"
	"ta/internal/domain"
	"ta/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user *domain.User) error {
	if user.Age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *UserService) GetByID(id string) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID)
}

