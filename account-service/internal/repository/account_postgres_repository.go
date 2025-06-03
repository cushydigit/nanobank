package repository

import (
	"context"
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type AccountRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.Account, error)
	Create(ctx context.Context, account *models.Account) error
	UpdateBalance(ctx context.Context, userID string, amount int64) error
}

type PostgresAccountRepository struct {
	DB *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{DB: db}
}

func (r *PostgresAccountRepository) FindByUserID(ctx context.Context, userID string) (*models.Account, error) {
	var a models.Account
	if err := r.DB.QueryRowContext(
		ctx,
		`SELECT id, user_id, username,  balance, created_at, updated_at FROM accounts WHERE user_id = $1`,
		userID,
	).Scan(&a.ID, &a.UserID, &a.Username, &a.Balance, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *PostgresAccountRepository) Create(ctx context.Context, a *models.Account) error {
	if _, err := r.DB.ExecContext(
		ctx,
		`INSERT INTO accounts (id, user_id, username, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		a.ID, a.UserID, a.Username, a.Balance, a.CreatedAt, a.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *PostgresAccountRepository) UpdateBalance(ctx context.Context, userID string, amount int64) error {
	if _, err := r.DB.ExecContext(
		ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE user_id = $2`,
		amount, userID,
	); err != nil {
		return err
	}
	return nil
}
