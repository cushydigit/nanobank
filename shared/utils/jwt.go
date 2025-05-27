package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateTokens(user *models.User) (*JWTTokens, error) {
	log.Printf("THE SECRET IS: %s", secret)
	now := time.Now()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 5)),
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
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
		},
	})
	refreshToken, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenStr, secret string) (*types.JWTClaims, error) {
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
