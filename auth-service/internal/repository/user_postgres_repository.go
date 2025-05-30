package repository

import (
	"context"
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// if user is not found it will return sql.ErrorNoRows
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.DB.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, created_at FROM users WHERE email = $1`,
		email,
	)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Passowrd, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, err := r.DB.ExecContext(
		ctx,
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
