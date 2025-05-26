package handler

import (
	"net/http"

	service "github.com/cushydigit/nanobank/auth-service/internal"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

}
