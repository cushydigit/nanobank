package repository

import (
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	return nil
}
