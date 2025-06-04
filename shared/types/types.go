package types

import (
	"github.com/cushydigit/nanobank/shared/models"
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
	RegisterReqKey          ContextKey = "register_req"
	AuthReqKey              ContextKey = "auth_req"
	RefreshReqKey           ContextKey = "refresh_req"
	UserIDKey               ContextKey = "user_id"
	UserEmailKey            ContextKey = "user_email"
	UpdateBlanceReqKey      ContextKey = "update_balance_req"
	InitiateTransferReqKey  ContextKey = "initiate_transfer_req"
	ConfirmTransferReqKey   ContextKey = "confirm_transfer_req"
	CreateTransactionReqKey ContextKey = "create_transaction_req"
	UpdateTransactionReqKey ContextKey = "update_transaction_req"
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

type InitiateTransferReqBody struct {
	Amount   int64  `json:"amount"`
	ToUserID string `json:"to_user_id"`
}

type ConfirmTransferReqBody struct {
	TxID  string `json:"tx_id"`
	Token string `json:"token"`
}

// internal requests
type CreateTransactionReqBody struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Amount     int64  `json:"amount"`
}

type UpdateTransactionReqBody struct {
	ID     string                   `json:"id"`
	Status models.TransactionStatus `json:"status"`
}

// response
type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// general response

// models types
