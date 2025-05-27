package handler

import (
	"errors"
	"net/http"

	"github.com/cushydigit/nanobank/auth-service/internal/service"
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
		helpers.ErrorJSON(w, errors.New("object not found in context of request"), http.StatusInternalServerError)
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
		Message: "registeration successfull",
		Data:    newUser,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.AuthReqKey).(types.AuthReqBody)
	if !ok {
		helpers.ErrorJSON(w, errors.New("object not found in context of request"), http.StatusInternalServerError)
		return
	}

	user, tokens, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case myerrors.ErrInternalServer:
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		default:
			helpers.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
	}

	payload := types.Response{
		Error:   false,
		Message: "login successfull",
		Data: map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"tokens":   tokens,
		},
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(types.RefreshReqKey).(types.RefreshReqBody)
	if !ok {
		helpers.ErrorJSON(w, errors.New("object not found in context of request"), http.StatusInternalServerError)
		return
	}
	user, tokens, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		switch err {
		case myerrors.ErrInternalServer:
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		default:
			helpers.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
	}

	payload := types.Response{
		Error:   false,
		Message: "refresh successfull",
		Data: map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"tokens":   tokens,
		},
	}

	helpers.WriteJSON(w, http.StatusOK, payload)

}
