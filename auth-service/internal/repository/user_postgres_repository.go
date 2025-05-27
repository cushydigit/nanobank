package repository

import (
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type UserRepository interface {
	// return err = nil if user is not found
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// return nil err and user if user not found
func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, username, email, password, created_at FROM users WHERE email = $1`,
		email,
	)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Passowrd, &user.CreatedAt)
	if err == sql.ErrNoRows {
		// the user not found
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	if _, err := r.DB.Exec(
		`INSERT INTO users (id, username, email, password, created_at) VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Username,
		user.Email,
		user.Passowrd,
		user.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
