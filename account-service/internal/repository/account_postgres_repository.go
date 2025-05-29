package repository

import (
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type AccountRepository interface {
	GetByUserID(userID string) (*models.Account, error)
	Create(*models.Account) error
	Deposit(amount float64) error
	Withdraw(amount float64) error
}

type PostgresAccountRepository struct {
	DB *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{DB: db}
}

func (r *PostgresAccountRepository) GetByUserID(userID string) (*models.Account, error) {
	return nil, nil
}

func (r *PostgresAccountRepository) Create(*models.Account) error {
	return nil
}

func (r *PostgresAccountRepository) Deposit(amount float64) error {
	return nil
}

func (r *PostgresAccountRepository) Withdraw(amount float64) error {
	return nil
}
