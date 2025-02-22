package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
)

var ErrExpenseCategoryWithThatNameAlreadyExists = errors.New("expense category with that name already exists")
var ErrUserIsNotVaultOwner = errors.New("user is not vault owner")

type ExpenseCategoryService struct {
	expenseCategoryRepo *repositories.ExpenseCategoryRepo
	vaultService        *VaultService
}

func NewExpenseCategoryService(expenseCategoryRepo *repositories.ExpenseCategoryRepo, vaultService *VaultService) *ExpenseCategoryService {
	return &ExpenseCategoryService{
		expenseCategoryRepo: expenseCategoryRepo,
		vaultService:        vaultService,
	}
}

func (s *ExpenseCategoryService) CreateOne(ctx context.Context, name string, userID, vaultID string) error {
	userVaultWithRole, err := s.vaultService.FindOneByID(ctx, userID, vaultID)
	if err != nil {
		return fmt.Errorf("failed to find vault %s for user %s: %w", vaultID, userID, err)
	}

	if userVaultWithRole.UserRole != models.VaultRoleOwner {
		return ErrUserIsNotVaultOwner
	}

	categories, err := s.expenseCategoryRepo.FindAll(ctx, userVaultWithRole.ID)
	if err != nil {
		return fmt.Errorf("failed to find existing expense categories in vault %s before creating one: %w", userVaultWithRole.ID, err)
	}

	exists := false
	for _, category := range categories {
		if category.Name == name {
			exists = true
			break
		}
	}
	if exists {
		return ErrExpenseCategoryWithThatNameAlreadyExists
	}

	err = s.expenseCategoryRepo.CreateOne(ctx, name, models.ExpenseCategoryStatusActive, 0, userVaultWithRole.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to create expense category: %w", err)
	}

	return nil
}

func (s *ExpenseCategoryService) FindAll(ctx context.Context, userID, vaultID string) ([]models.ExpenseCategory, error) {
	vault, err := s.vaultService.FindOneByID(ctx, userID, vaultID)
	if err != nil {
		return nil, err
	}

	categories, err := s.expenseCategoryRepo.FindAll(ctx, vault.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find expense categories for vault %s & user %s: %w", vault.ID, userID, err)
	}

	return categories, nil
}
