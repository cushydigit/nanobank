package service

import (
	"context"

	"github.com/cushydigit/nanobank/account-service/internal/repository"
	"github.com/cushydigit/nanobank/shared/models"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(r repository.AccountRepository) *AccountService {
	return &AccountService{repo: r}
}

func GetByUserID(ctx context.Context, userID string) (*models.Account, error) {
	return nil, nil
}

func Create(ctx context.Context, userID string) (*models.Account, error) {
	return nil, nil
}

func Deposit(ctx context.Context, amount string) error {
	return nil
}

func Withdraw(ctx context.Context, amount string) error {
	return nil
}
