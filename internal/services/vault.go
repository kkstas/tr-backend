package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
)

var ErrVaultNotFound = errors.New("vault not found")
var ErrInsufficientVaultPermissions = errors.New("insufficient permissions to perform this vault operation")
var ErrUserAlreadyAssignedToVault = errors.New("user is already assigned to this vault")

type VaultService struct {
	vaultRepo   *repositories.VaultRepo
	userService *UserService
}

func NewVaultService(vaultRepo *repositories.VaultRepo, userService *UserService) *VaultService {
	return &VaultService{vaultRepo: vaultRepo, userService: userService}
}

func (s *VaultService) CreateOne(ctx context.Context, userID, vaultName string) error {
	user, err := s.userService.FindOneByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user %s before creating vault: %w", userID, err)
	}

	vaultID, err := s.vaultRepo.CreateOne(ctx, userID, models.VaultRoleOwner, vaultName)
	if err != nil {
		return fmt.Errorf("failed to create new vault: %w", err)
	}

	if user.ActiveVault == "" {
		err := s.userService.AssignActiveVault(ctx, userID, vaultID)
		if err != nil {
			return fmt.Errorf("failed to assign active vault %s to user %s after creating new vault: %w", vaultID, userID, err)
		}
	}

	return nil
}

func (s *VaultService) FindAll(ctx context.Context, userID string) ([]models.UserVaultWithRole, error) {
	return s.vaultRepo.FindAll(ctx, userID)
}

func (s *VaultService) FindOneByID(ctx context.Context, userID, vaultID string) (*models.UserVaultWithRole, error) {
	vault, err := s.vaultRepo.FindOneByID(ctx, userID, vaultID)
	if err != nil {
		if errors.Is(err, repositories.ErrVaultNotFound) {
			return nil, ErrVaultNotFound
		}
		return nil, fmt.Errorf("failed to find vault %s for user %s: %w", vaultID, userID, err)
	}
	return vault, nil
}

func (s *VaultService) DeleteOneByID(ctx context.Context, userID, vaultID string) error {
	foundVault, err := s.vaultRepo.FindOneByID(ctx, userID, vaultID)
	if err != nil {
		if errors.Is(err, repositories.ErrVaultNotFound) {
			return ErrVaultNotFound
		}
		return err
	}
	if foundVault.UserRole != models.VaultRoleOwner {
		return ErrInsufficientVaultPermissions
	}

	err = s.vaultRepo.DeleteOneByID(ctx, vaultID)
	if err != nil {
		return fmt.Errorf("failed to delete vault %s as user %s: %w", vaultID, userID, err)
	}
	return nil

}

func (s *VaultService) AddUser(ctx context.Context, userID, invitedUserID, vaultID string, userRole models.VaultRole) error {
	_, err := s.vaultRepo.FindOneByID(ctx, invitedUserID, vaultID)
	if err == nil {
		return ErrUserAlreadyAssignedToVault
	}

	userVaultWithRole, err := s.vaultRepo.FindOneByID(ctx, userID, vaultID)
	if err != nil {
		if errors.Is(err, repositories.ErrVaultNotFound) {
			return ErrVaultNotFound
		}
		return fmt.Errorf("failed to find user vault: %w", err)
	}

	if userVaultWithRole.UserRole != models.VaultRoleOwner {
		return ErrInsufficientVaultPermissions
	}

	return s.vaultRepo.AddUser(ctx, userVaultWithRole.ID, invitedUserID, userRole)
}
