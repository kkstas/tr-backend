package repositories

import (
	"database/sql"
	"fmt"

	"github.com/kkstas/tnr-backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateOne(firstName, lastName, email string) error {
	_, err := u.db.Exec(`INSERT INTO users(first_name, last_name, email) VALUES ($1, $2, $3);`, firstName, lastName, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepo) FindAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, first_name, last_name, email FROM users;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email); err != nil {
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
