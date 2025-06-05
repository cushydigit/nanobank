package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/redis"
	"github.com/cushydigit/nanobank/shared/utils"
	"github.com/cushydigit/nanobank/transaction-service/internal/repository"
)

type TransactionService struct {
	repo        repository.TransactionRepository
	tokenCacher redis.TokenCacher
}

func NewTransactionService(r repository.TransactionRepository, c redis.TokenCacher) *TransactionService {
	return &TransactionService{repo: r}
}

// returns ErrTransactionNotFound, ErrInternalServer
func (s *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myerrors.ErrTransactionNotFound
		}
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}

	return t, nil
}

// returns ErrInternalServer
func (s *TransactionService) ListAll(ctx context.Context) ([]*models.Transaction, error) {
	ts, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return ts, nil
}

// reutrns ErrInternalServer
func (s *TransactionService) ListByUserID(ctx context.Context, userID string) ([]*models.Transaction, error) {
	ts, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	var filteredTs []*models.Transaction
	for _, t := range ts {
		if t.FromUserID == userID {
			filteredTs = append(filteredTs, t)
		} else {
			if t.Status == models.StatusConfirmed {
				filteredTs = append(filteredTs, t)
			}
		}
	}
	return filteredTs, nil
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

// returns ErrTransactionNotFound, ErrInternalServer
func (s *TransactionService) Update(ctx context.Context, id string, status models.TransactionStatus) (*models.Transaction, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myerrors.ErrTransactionNotFound
		}
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	t.UpdatedAt = time.Now().UTC()
	t.Status = status
	if err := s.repo.Update(ctx, t); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	return t, nil
}
