package handler

import (
	"net/http"

	"github.com/cushydigit/nanobank/account-service/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

func (h *AccountHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {

}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {

}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {

}
