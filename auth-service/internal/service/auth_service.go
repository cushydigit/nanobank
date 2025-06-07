package service

import (
	"context"
	"database/sql"
	"log"

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

// returns ErrDuplicateEmail, ErrInternalServer
func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {

	// check email duplication
	if exists, _ := s.repo.FindByEmail(ctx, email); exists != nil {
		return nil, myerrors.ErrDuplicateEmail
	}

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}
	newUser := models.NewUser(username, email, hashedPassword)

	// insert new user to DB
	if err := s.repo.Create(ctx, newUser); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, myerrors.ErrInternalServer
	}

	return newUser, nil

}

// returns ErrInvalidCredentials, ErrInternalServer
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, *types.JWTTokens, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, myerrors.ErrInvalidCredentials
		} else {
			log.Printf("unexpected err: %v", err)
			return nil, nil, myerrors.ErrInternalServer
		}
	}
	// check password
	if ok := utils.CheckPasswordHash(password, user.Passowrd); !ok {
		return nil, nil, myerrors.ErrInvalidCredentials
	}

	// generate tokens
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	// store new tokens in cache
	if err := s.cacher.SetAuth(ctx, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	return user, tokens, nil

}

// returns ErrInvalidRefreshToken, ErrInternalServer
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.User, *types.JWTTokens, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, nil, myerrors.ErrInvalidRefreshToken
	}

	// check if token is not rotated
	token, err := s.cacher.GetAuth(ctx, claims.UserID)
	if err != nil || token != refreshToken {
		return nil, nil, myerrors.ErrInvalidRefreshToken
	}

	user := &models.User{
		ID:       claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
	}
	// genereate new tokens
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	// store new tokens in cache
	if err := s.cacher.SetAuth(ctx, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("unexpected err: %v", err)
		return nil, nil, myerrors.ErrInternalServer
	}

	return user, tokens, nil
}

// returns ErrInvalidRefreshToken, ErrInternalServer
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return myerrors.ErrInvalidRefreshToken
	}

	// delete tokens from cache
	if err := s.cacher.DelAuth(ctx, claims.UserID); err != nil {
		log.Printf("unexpected err: %v", err)
		return myerrors.ErrInternalServer
	}

	return nil
}
