package handler

import (
	"net/http"

	"github.com/cushydigit/nanobank/account-service/internal/service"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	// will injected in gateway with requireauth middleware
	userID := r.Header.Get(string(types.XUserID))

	a, err := h.service.Get(r.Context(), userID)
	if err != nil {
		if err == myerrors.ErrAccountNotFound {
			helpers.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := &types.Response{
		Error:   false,
		Message: "success",
		Data:    a,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// will injected in gateway with requireauth middleware
	userID := r.Header.Get(string(types.XUserID))

	a, err := h.service.Create(r.Context(), userID)
	if err != nil {
		if err == myerrors.ErrAccountAlreadyExists {
			helpers.ErrorJSON(w, err, http.StatusConflict)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := &types.Response{
		Error:   false,
		Message: "success",
		Data:    a,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	// will injected in gateway with requireauth middleware
	userID := r.Header.Get(string(types.XUserID))
	req, ok := r.Context().Value(string(types.UpdateBlanceReqKey)).(types.UpdateBalanceReqBody)
	if !ok {
		helpers.ErrorJSON(w, myerrors.ErrContextValueNotFoundInRequest, http.StatusInternalServerError)
		return
	}

	if err := h.service.Deposit(r.Context(), userID, req.Amount); err != nil {
		if err == myerrors.ErrAccountNotFound {
			helpers.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		if err == myerrors.ErrAmountMustBePositive {
			helpers.ErrorJSON(w, err)
			return
		}
		helpers.ErrorJSON(w, myerrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    nil,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	// will injected in gateway with requireauth middleware
	userID := r.Header.Get(string(types.XUserID))
	req, ok := r.Context().Value(string(types.UpdateBlanceReqKey)).(types.UpdateBalanceReqBody)
	if !ok {
		helpers.ErrorJSON(w, myerrors.ErrContextValueNotFoundInRequest, http.StatusInternalServerError)
		return
	}
	if err := h.service.Withdraw(r.Context(), userID, req.Amount); err != nil {
		if err == myerrors.ErrAccountNotFound {
			helpers.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		if err == myerrors.ErrAmountMustBePositive {
			helpers.ErrorJSON(w, err)
			return
		}
		if err == myerrors.ErrInsufficientBalance {
			helpers.ErrorJSON(w, err, http.StatusUnprocessableEntity)
			return
		}
		helpers.ErrorJSON(w, myerrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    nil,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}
