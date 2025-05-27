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

type HeaderKey string

const (
	XUserID    HeaderKey = "X-User-ID"
	XUserEmail HeaderKey = "X-Email-ID"
)

type ContextKey string

const (
	RegisterReqKey ContextKey = "register_req"
	AuthReqKey     ContextKey = "auth_req"
	UserIDKey      ContextKey = "user_id"
	UserEmailKey   ContextKey = "user_email"
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

// response
type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// general response

// models types
