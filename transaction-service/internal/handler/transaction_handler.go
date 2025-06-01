package handler

import (
	"errors"
	"net/http"

	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/transaction-service/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (t *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	helpers.ErrorJSON(w, errors.New("not implemenetd yet"))
	return
}

func (t *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	helpers.ErrorJSON(w, errors.New("not implemented yet"))
	return
}
