package handler

import (
	"errors"
	"net/http"

	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/cushydigit/nanobank/transaction-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.CreateTransactionReqKey).(types.CreateTransactionReqBody)
	if !ok {
		helpers.ErrorJSON(w, myerrors.ErrContextValueNotFoundInRequest, http.StatusInternalServerError)
		return
	}
	t, err := h.service.Create(r.Context(), req.FromUserID, req.ToUserID, req.Amount)
	if err != nil {
		if err == myerrors.ErrAmountMustBePositive {
			helpers.ErrorJSON(w, err)
			return
		} else {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		}
	}
	// TODO add the transaction with ttl

	payload := types.Response{
		Error:   false,
		Message: "trnasction created",
		Data:    t,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.UpdateTransactionReqKey).(types.UpdateTransactionReqBody)
	if !ok {
		helpers.ErrorJSON(w, myerrors.ErrContextValueNotFoundInRequest, http.StatusInternalServerError)
		return
	}

	t, err := h.service.Update(r.Context(), req.ID, req.Status)
	if err != nil {
		if err == myerrors.ErrTransactionNotFound {
			helpers.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    t,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	t, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == myerrors.ErrTransactionNotFound {
			helpers.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    t,
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *TransactionHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {

}

func (t *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	helpers.ErrorJSON(w, errors.New("not implemented yet"))
	return
}
