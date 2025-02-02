package services

import (
	"context"
	"fmt"

	"github.com/kkstas/tnr-backend/internal/auth"
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

func (s *UserService) CreateOne(ctx context.Context, firstName, lastName, email, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.userRepo.CreateOne(ctx, firstName, lastName, email, passwordHash)
}

func (u *UserService) FindPasswordHashAndUserIDForEmail(ctx context.Context, email string) (passwordHash, userID string, err error) {
	passwordHash, userID, err = u.userRepo.FindPasswordHashAndUserIDForEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to find password hash: %w", err)
	}
	return passwordHash, userID, nil
}

func (s *UserService) FindOneByID(ctx context.Context, id string) (models.User, error) {
	return s.userRepo.FindOneByID(ctx, id)
}
