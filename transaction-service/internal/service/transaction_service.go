package service

import (
	"context"
	"log"
	"time"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/utils"
	"github.com/cushydigit/nanobank/transaction-service/internal/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(r repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: r}
}

func (s *TransactionService) ListAll(ctx context.Context) ([]*models.Transaction, error) {
	ts, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return ts, nil
}

func (s *TransactionService) ListByUserID(ctx context.Context, userID string) ([]*models.Transaction, error) {
	ts, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return ts, nil
}

// returns ErrAmountMustBePositive, ErrInternalServer
func (s *TransactionService) Create(ctx context.Context, fromUserID, toUserID string, amount int64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, myerrors.ErrAmountMustBePositive
	}
	token, err := utils.GenerateTransactionToken()
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	t := models.NewTransaction(fromUserID, toUserID, token, amount)
	if err := s.repo.Create(ctx, t); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return t, nil
}

func (s *TransactionService) Update(ctx context.Context, id string, status models.TransactionStatus) error {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return myerrors.ErrTransactionNotFound
	}
	t.UpdatedAt = time.Now().UTC()
	t.Status = status
	if err := s.repo.Update(ctx, t); err != nil {
		log.Printf("unexpected err: %v")
		return myerrors.ErrTransactionNotFound
	}
	return nil
}
