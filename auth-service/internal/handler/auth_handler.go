package handler

import (
	"errors"
	"log"
	"net/http"

	service "github.com/cushydigit/nanobank/auth-service/internal"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/helpers"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/cushydigit/nanobank/shared/utils"
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
		helpers.ErrorJSON(w, errors.New("request object not found in context"), http.StatusInternalServerError)
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
		helpers.ErrorJSON(w, errors.New("request object not found in context"), http.StatusInternalServerError)
		return
	}

	user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if err == myerrors.ErrInternalServer {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
	}

	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		log.Printf("error is: %v", err)
		helpers.ErrorJSON(w, errors.New("could not generate tokens"), http.StatusInternalServerError)
		return
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

}
