package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/cushydigit/nanobank/account-service/internal/repository"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(r repository.AccountRepository) *AccountService {
	return &AccountService{repo: r}
}

func (s *AccountService) Get(ctx context.Context, userID string) (*models.Account, error) {
	a, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		// the account not found
		if err == sql.ErrNoRows {
			return nil, myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return nil, err
		}
	}
	return a, nil
}

func (s *AccountService) Create(ctx context.Context, userID string) (*models.Account, error) {
	if exists, _ := s.repo.FindByUserID(ctx, userID); exists != nil {
		return nil, myerrors.ErrAccountAlreadyExists
	}
	newAccount := models.NewAccount(userID)
	if err := s.repo.Create(ctx, newAccount); err != nil {
		log.Printf("unexpected error: %v", err)
		return nil, err
	}
	return newAccount, nil
}

func (s *AccountService) Deposit(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return myerrors.ErrAmountMustBePositive
	}
	if _, err := s.repo.FindByUserID(ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return err
		}
	}
	if err := s.repo.UpdateBalance(ctx, userID, amount); err != nil {
		log.Printf("unexpected error: %v", err)
		return err
	}
	return nil
}

func (s *AccountService) Withdraw(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return myerrors.ErrAmountMustBePositive
	}
	a, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return myerrors.ErrAccountNotFound
		} else {
			log.Printf("unexpected error: %v", err)
			return err
		}
	}

	if a.Balance < amount {
		return myerrors.ErrInsufficientBalance
	}
	if err := s.repo.UpdateBalance(ctx, userID, -amount); err != nil {
		log.Printf("unexpected error: %v", err)
		return err
	}
	return nil
}
