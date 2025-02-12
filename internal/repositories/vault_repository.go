package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
)

type VaultRepo struct {
	db *sql.DB
}

func NewVaultRepo(db *sql.DB) *VaultRepo {
	return &VaultRepo{db: db}
}

func (r *VaultRepo) CreateOne(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO vaults(id, name) VALUES ($1, $2);`, uuid.New(), name)
	if err != nil {
		// TODO: handle unique constraint errors
		return fmt.Errorf("failed to create vault: %w", err)
	}
	return nil
}

func (r *VaultRepo) FindAll(ctx context.Context) ([]models.Vault, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at FROM vaults;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query vaults: %w", err)
	}

	defer rows.Close()

	vaults := []models.Vault{}

	for rows.Next() {
		var v models.Vault
		if err := rows.Scan(&v.ID, &v.Name, &v.CreatedAt); err != nil {
			return nil, err
		}
		vaults = append(vaults, v)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return vaults, nil
}
