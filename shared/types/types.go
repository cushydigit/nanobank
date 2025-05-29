package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type HeaderKey string

const (
	XUserID    HeaderKey = "X-User-ID"
	XUserEmail HeaderKey = "X-Email-ID"
	XUsername  HeaderKey = "X-Username"
)

type ContextKey string

const (
	RegisterReqKey     ContextKey = "register_req"
	AuthReqKey         ContextKey = "auth_req"
	RefreshReqKey      ContextKey = "refresh_req"
	UserIDKey          ContextKey = "user_id"
	UserEmailKey       ContextKey = "user_email"
	UpdateBlanceReqKey ContextKey = "update_balance_req"
)

// request bodies
type RegisterReqBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshReqBody struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateBalanceReqBody struct {
	Amount int64 `json:"amount"`
}

// response
type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// general response

// models types
