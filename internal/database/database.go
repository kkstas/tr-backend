package database

import (
	"context"
	"database/sql"
	"fmt"
)

func InitDBTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id         TEXT PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name  VARCHAR(255) NOT NULL,
			email      VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS vaults (
			id          TEXT PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create vaults table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS user_vaults (
			user_id     TEXT NOT NULL,
			vault_id    TEXT NOT NULL,
			role        VARCHAR(255) NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, vault_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create vaults table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS expense_categories (
			id         TEXT PRIMARY KEY,
			name       VARCHAR(255) NOT NULL UNIQUE,
			priority   INTEGER NOT NULL,
			vault_id   TEXT NOT NULL,
			created_by TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create expense_categories table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS expenses (
			id             TEXT PRIMARY KEY,
			name           TEXT NOT NULL,
			date           TEXT NOT NULL,
			category_id    TEXT NOT NULL,
			amount         REAL NOT NULL,
			payment_method TEXT NOT NULL,
			vault_id       TEXT NOT NULL,
			created_by     TEXT NOT NULL,
			created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES users(id),
			FOREIGN KEY (category_id) REFERENCES expense_categories(id)
		);
	`)

	if err != nil {
		return fmt.Errorf("failed to create expenses table: %w", err)
	}

	return nil
}
