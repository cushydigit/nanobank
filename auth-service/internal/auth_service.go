package service

import "github.com/cushydigit/nanobank/auth-service/internal/repository"

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(r repository.UserRepository) *AuthService {
	return &AuthService{repo: r}
}
