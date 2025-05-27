package utils

import (
	"errors"
	"os"
	"time"

	"github.com/cushydigit/nanobank/shared/config"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

func GenerateTokens(user *models.User) (*types.JWTTokens, error) {
	now := time.Now()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.TTL_ACCESS_TOKEN)),
		},
	})
	accessToken, err := access.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.TTL_REFRESH_TOKEN)),
		},
	})
	refreshToken, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &types.JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenStr string) (*types.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
