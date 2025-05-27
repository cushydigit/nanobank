package service

import (
	"github.com/cushydigit/nanobank/auth-service/internal/repository"
	"github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/utils"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(r repository.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {

	// check email duplication
	if exists, _ := s.repo.FindByEmail(email); exists != nil {
		return nil, errors.ErrDuplicateEmail
	}

	// hash password
	hashedPassword, _ := utils.HashPassword(password)
	newUser := models.NewUser(username, email, hashedPassword)

	// insert new user to DB
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil

}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if user == nil {
		// user not found
		return nil, errors.ErrInvalidCredentials
	}
	// check password
	if ok := utils.CheckPasswordHash(password, user.Passowrd); !ok {
		return nil, errors.ErrInvalidCredentials
	}

	return user, nil

}
