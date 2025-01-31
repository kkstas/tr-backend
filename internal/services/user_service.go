package services

import (
	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepo
}

func NewUserService(userRepo *repositories.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) FindAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) CreateOne(firstName string, lastName string, email string) error {
	return s.userRepo.CreateOne(firstName, lastName, email)
}
