package services

import (
	"context"

	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepo
}

func NewUserService(userRepo *repositories.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) FindAll(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *UserService) CreateOne(ctx context.Context, firstName string, lastName string, email string) error {
	return s.userRepo.CreateOne(ctx, firstName, lastName, email)
}

func (s *UserService) FindOneByID(ctx context.Context, id string) (models.User, error) {
	return s.userRepo.FindOneByID(ctx, id)
}
