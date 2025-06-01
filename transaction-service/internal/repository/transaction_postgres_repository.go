package repository

import (
	"context"
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type TransactionRepository interface {
	FindByID(ctx context.Context, id string) (*models.Transaction, error)
}

type PostgresTransactionRepository struct {
	DB *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{DB: db}
}

func (r *PostgresTransactionRepository) FindByID(ctx context.Context, id string) (*models.Transaction, error) {
	return nil, nil
}
