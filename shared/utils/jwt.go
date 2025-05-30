package utils

import (
	"os"
	"time"

	"github.com/cushydigit/nanobank/shared/config"
	myerrors "github.com/cushydigit/nanobank/shared/errors"
	"github.com/cushydigit/nanobank/shared/models"
	"github.com/cushydigit/nanobank/shared/types"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func GenerateTokens(user *models.User) (*types.JWTTokens, error) {
	now := time.Now()

	secret := []byte(JWT_SECRET)

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.TTL_ACCESS_TOKEN)),
		},
	})
	accessToken, err := access.SignedString(secret)
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
	refreshToken, err := refresh.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &types.JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// returns ErrJWTEmptyToken, ErrJWTExpiredToken, ErrJWTInvalidToken, ErrJWTTokenClaimsTypeMismatch, ErrJWTFailedToParseToken
func ValidateToken(tokenStr string) (*types.JWTClaims, error) {

	if tokenStr == "" {
		return nil, myerrors.ErrJWTEmptyToken
	}

	secret := []byte(JWT_SECRET)

	token, err := jwt.ParseWithClaims(tokenStr, &types.JWTClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, myerrors.ErrJWTFailedToParseToken
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok || !token.Valid {
		return nil, myerrors.ErrJWTTokenClaimsTypeMismatch
	}

	if !token.Valid {
		return nil, myerrors.ErrJWTInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, myerrors.ErrJWTExpiredToken
	}

	return claims, nil
}
