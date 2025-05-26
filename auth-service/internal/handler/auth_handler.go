package handler

import (
	"errors"
	"net/http"

	service "github.com/cushydigit/nanobank/auth-service/internal"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.RegisterReqKey).(types.RegisterReqBody)
	if !ok {
		helpers.ErrorJSON(w, errors.New("registeration request not found in context"), http.StatusInternalServerError)
		return
	}

	newUser, err := h.service.Register(req.Username, req.Email, req.Password)
	if err == myerrors.ErrDuplicateEmail {
		helpers.ErrorJSON(w, err, http.StatusConflict)
		return
	} else if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "success",
		Data:    newUser,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

}
