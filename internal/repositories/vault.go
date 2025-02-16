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

func (r *VaultRepo) CreateOne(ctx context.Context, userID string, userRole models.VaultRole, vaultName string) (vaultID string, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create vault: %w", err)
	}

	defer tx.Rollback()

	vaultID = uuid.New().String()

	_, err = tx.ExecContext(ctx, `INSERT INTO vaults(id, name) VALUES ($1, $2)`, vaultID, vaultName)
	if err != nil {
		return "", fmt.Errorf("failed to insert new vault: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_vaults(user_id, vault_id, role) VALUES ($1, $2, $3)`, userID, vaultID, userRole)
	if err != nil {
		return "", fmt.Errorf("failed to insert new record in user_vaults junction table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return vaultID, nil
}

func (r *VaultRepo) FindAll(ctx context.Context, userID string) ([]models.UserVaultWithRole, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT vaults.id, vaults.name, user_vaults.role FROM users
		JOIN user_vaults ON users.id = user_vaults.user_id
		JOIN vaults ON vaults.id = user_vaults.vault_id
		WHERE users.id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query vaults: %w", err)
	}

	defer rows.Close()

	vaults := []models.UserVaultWithRole{}

	for rows.Next() {
		var v models.UserVaultWithRole
		if err := rows.Scan(&v.ID, &v.Name, &v.UserRole); err != nil {
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

func (r *VaultRepo) FindOneByID(ctx context.Context, userID, vaultID string) (*models.UserVaultWithRole, error) {
	v := models.UserVaultWithRole{}

	err := r.db.QueryRowContext(ctx, `
		SELECT v.id, v.name, uv.role FROM vaults v
		JOIN user_vaults uv ON uv.vault_id = v.id
		WHERE v.id = $1 AND uv.user_id = $2
		`, vaultID, userID).Scan(&v.ID, &v.Name, &v.UserRole)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *VaultRepo) DeleteOneByID(ctx context.Context, vaultID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM vaults WHERE vaults.id = $1`, vaultID)
	return err
}
