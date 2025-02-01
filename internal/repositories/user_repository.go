package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/kkstas/tnr-backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateOne(ctx context.Context, firstName, lastName, email string) error {
	_, err := u.db.ExecContext(ctx, `
		INSERT INTO users(id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4);`,
		uuid.New().String(), firstName, lastName, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepo) FindAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, first_name, last_name, email, created_at
		FROM users;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) FindOne(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
			SELECT id, first_name, last_name, email, created_at
			FROM users
			WHERE users.id = $1;
		`, id).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
