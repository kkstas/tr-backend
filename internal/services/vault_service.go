package services

import (
	"context"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
)

type VaultService struct {
	vaultRepo *repositories.VaultRepo
}

func NewVaultService(vaultRepo *repositories.VaultRepo) *VaultService {
	return &VaultService{vaultRepo: vaultRepo}
}

func (s *VaultService) FindAll(ctx context.Context) ([]models.Vault, error) {
	return s.vaultRepo.FindAll(ctx)
}

func (s *VaultService) CreateOne(ctx context.Context, name string) error {
	return s.vaultRepo.CreateOne(ctx, name)
}
