package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
)

var ErrExpenseCategoryNotFound = errors.New("expense category not found")

type ExpenseCategoryRepo struct {
	db *sql.DB
}

func NewExpenseCategoryRepo(db *sql.DB) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{db: db}
}

func (r *ExpenseCategoryRepo) CreateOne(ctx context.Context, name string, status models.ExpenseCategoryStatus, priority int, vaultID, createdBy string) (categoryID string, err error) {
	categoryID = uuid.New().String()

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO expense_categories(id, name, status, priority, vault_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		categoryID, name, status, priority, vaultID, createdBy)
	if err != nil {
		return "", err
	}
	return categoryID, nil
}

func (r *ExpenseCategoryRepo) FindAll(ctx context.Context, vaultID string) ([]models.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, status, priority, vault_id, created_by, created_at
		FROM expense_categories
		WHERE vault_id = $1`, vaultID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find all categories query for vault %s: %w", vaultID, err)
	}
	defer rows.Close()

	categories := []models.ExpenseCategory{}

	var category = models.ExpenseCategory{} // nolint: exhaustruct
	for rows.Next() {
		err := rows.Scan(&category.ID, &category.Name, &category.Status, &category.Priority, &category.VaultID, &category.CreatedBy, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ExpenseCategoryRepo) FindOneByID(ctx context.Context, categoryID string) (*models.ExpenseCategory, error) {
	category := models.ExpenseCategory{} // nolint: exhaustruct

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, status, priority, vault_id, created_by, created_at
		FROM expense_categories
		WHERE id = $1
		`, categoryID).Scan(&category.ID, &category.Name, &category.Status, &category.Priority, &category.VaultID, &category.CreatedBy, &category.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrExpenseCategoryNotFound
		}

		return nil, err
	}

	return &category, nil
}

func (r *ExpenseCategoryRepo) SetStatus(ctx context.Context, categoryID string, status models.ExpenseCategoryStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE expense_categories
		SET status = $1
		WHERE id = $2`, status, categoryID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExpenseCategoryRepo) SetPriority(ctx context.Context, categoryID string, priority int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE expense_categories
		SET priority = $1
		WHERE id = $2`, priority, categoryID,
	)
	if err != nil {
		return err
	}
	return nil
}
