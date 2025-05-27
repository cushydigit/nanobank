package service

import (
	"context"

	"github.com/cushydigit/nanobank/auth-service/internal/repository"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	myredis "github.com/cushydigit/nanobank/shared/redis"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/cushydigit/nanobank/shared/utils"
)

type AuthService struct {
	repo   repository.UserRepository
	cacher myredis.AuthCacher
}

func NewAuthService(r repository.UserRepository, c myredis.AuthCacher) *AuthService {
	return &AuthService{
		repo:   r,
		cacher: c,
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {

	// check email duplication
	if exists, _ := s.repo.FindByEmail(email); exists != nil {
		return nil, myerrors.ErrDuplicateEmail
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

func (s *AuthService) Login(email, password string) (*models.User, *types.JWTTokens, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, nil, myerrors.ErrInternalServer
	}
	if user == nil {
		// user not found
		return nil, nil, myerrors.ErrInvalidCredentials
	}
	// check password
	if ok := utils.CheckPasswordHash(password, user.Passowrd); !ok {
		return nil, nil, myerrors.ErrInvalidCredentials
	}

	// generate tokens
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		return nil, nil, myerrors.ErrInternalServer
	}

	return user, tokens, nil

}

func (s *AuthService) Refresh(refreshToken string) (*models.User, *types.JWTTokens, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, nil, myerrors.ErrInvalidRefreshToken
	}

	user := &models.User{
		ID:       claims.ID,
		Email:    claims.Email,
		Username: claims.Username,
	}
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		return nil, nil, myerrors.ErrInternalServer
	}

	return user, tokens, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return myerrors.ErrInvalidRefreshToken
	}

	// delete from redis
	return s.cacher.DelAuth(context.Background(), claims.UserID)

}
