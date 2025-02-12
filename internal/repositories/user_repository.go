package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/kkstas/tr-backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailAlreadyExists = errors.New("user with that email already exists")

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateOne(ctx context.Context, firstName, lastName, email, passwordHash string) error {
	_, err := u.FindOneByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return fmt.Errorf("failed to find user before creating one: %w", err)
	}
	if err == nil {
		return ErrUserEmailAlreadyExists
	}

	_, err = u.db.ExecContext(ctx, `
		INSERT INTO users(id, first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5);`,
		uuid.New().String(), firstName, lastName, email, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u *UserRepo) FindPasswordHashAndUserIDForEmail(ctx context.Context, email string) (passwordHash, userID string, err error) {
	err = u.db.QueryRowContext(ctx, `SELECT u.id, u.password_hash FROM users u WHERE u.email = $1;`, email).Scan(&userID, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to find user password hash: %w", err)
	}
	return passwordHash, userID, nil
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
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt); err != nil {
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

func (r *UserRepo) FindOneByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
			SELECT id, first_name, last_name, email, created_at
			FROM users
			WHERE users.id = $1;
		`, id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
			SELECT id, first_name, last_name, email, created_at
			FROM users
			WHERE users.email = $1;
		`, email).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
