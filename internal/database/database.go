package database

import (
	"context"
	"database/sql"
	"fmt"
)

func OpenDB(ctx context.Context, dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbname+"?_pragma=foreign_keys(1)&_time_format=sqlite")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = initDBTables(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to init db tables: %w", err)
	}

	return db, nil
}

func initDBTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id            TEXT PRIMARY KEY,
			first_name    TEXT NOT NULL,
			last_name     TEXT NOT NULL,
			email         TEXT NOT NULL UNIQUE,
			active_vault  TEXT NULL,
			password_hash TEXT NOT NULL,
			created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (active_vault) REFERENCES vaults(id) ON DELETE SET NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS vaults (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create vaults table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS user_vaults (
			user_id     TEXT NOT NULL,
			vault_id    TEXT NOT NULL,
			role        TEXT NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, vault_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create vaults table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS expense_categories (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			status     TEXT DEFAULT 'active',
			priority   INTEGER NOT NULL,
			vault_id   TEXT NOT NULL,
			created_by TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES users(id)
		)
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
		)
	`)

	if err != nil {
		return fmt.Errorf("failed to create expenses table: %w", err)
	}

	return nil
}
