package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkstas/tr-backend/internal/auth"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailAlreadyExists = errors.New("user with that email already exists")

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

	_, err = s.userRepo.FindOneByEmail(ctx, email)
	if err == nil {
		return ErrUserEmailAlreadyExists
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return fmt.Errorf("failed to find user before creating one: %w", err)
	}

	return s.userRepo.CreateOne(ctx, firstName, lastName, email, passwordHash)
}

func (s *UserService) FindPasswordHashAndUserIDForEmail(ctx context.Context, email string) (passwordHash, userID string, err error) {
	passwordHash, userID, err = s.userRepo.FindPasswordHashAndUserIDForEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to find password hash: %w", err)
	}
	return passwordHash, userID, nil
}

func (s *UserService) FindOneByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.FindOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by id %s: %w", id, err)
	}

	return user, nil
}

func (s *UserService) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email %s: %w", email, err)
	}

	return user, nil
}

func (s *UserService) AssignActiveVault(ctx context.Context, userID, vaultID string) error {
	return s.userRepo.AssignActiveVault(ctx, userID, vaultID)
}
