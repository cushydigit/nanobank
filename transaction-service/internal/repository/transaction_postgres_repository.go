package repository

import (
	"context"
	"database/sql"

	"github.com/cushydigit/nanobank/shared/models"
)

type TransactionRepository interface {
	FindAll(ctx context.Context) ([]*models.Transaction, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*models.Transaction, error)
	FindByID(ctx context.Context, id string) (*models.Transaction, error)
	Create(ctx context.Context, t *models.Transaction) error
	Update(ctx context.Context, t *models.Transaction) error
}

type PostgresTransactionRepository struct {
	DB *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{DB: db}
}

func (r *PostgresTransactionRepository) FindByID(ctx context.Context, id string) (*models.Transaction, error) {
	var t models.Transaction
	if err := r.DB.QueryRowContext(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, status, confirmation_token, created_at, updated_at From transactions WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.ConfirmationToken, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresTransactionRepository) FindAll(ctx context.Context) ([]*models.Transaction, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, status, confirmation_token, created_at, updated_at FROM transactions`,
	)

	if err != nil {
		return nil, err
	}

	var ts []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.ConfirmationToken, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (r *PostgresTransactionRepository) FindAllByUserID(ctx context.Context, userID string) ([]*models.Transaction, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, status, confirmation_token, created_at, updated_at FROM transactions WHERE from_user_id = $1 OR to_user_id = $2`,
		userID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	var ts []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &t.ConfirmationToken, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, t *models.Transaction) error {
	if _, err := r.DB.ExecContext(
		ctx,
		`INSERT INTO transactions (id , from_user_id, to_user_id, amount, status, confirmation_token, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		t.ID, t.FromUserID, t.ToUserID, t.Amount, t.Status, t.ConfirmationToken, t.CreatedAt, t.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *PostgresTransactionRepository) Update(ctx context.Context, t *models.Transaction) error {
	if _, err := r.DB.ExecContext(
		ctx,
		`UPDATE transactions SET status = $1, updated_at= $2 WHERE id = $3`,
		t.Status,
		t.UpdatedAt,
		t.ID,
	); err != nil {
		return err
	}
	return nil
}
