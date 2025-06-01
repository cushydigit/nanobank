package service

import "github.com/cushydigit/nanobank/transaction-service/internal/repository"

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(r repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: r}
}
