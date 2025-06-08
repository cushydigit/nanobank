package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const CREATE_TRANSACTIONS_TABLE = `
		CREATE TABLE IF NOT EXISTS transactions (
  		id TEXT PRIMARY KEY,
  		from_user_id TEXT NOT NULL,
  		to_user_id TEXT NOT NULL,
  		amount BIGINT NOT NULL,
  		status transaction_status NOT NULL,
  		confirmation_token UUID NOT NULL UNIQUE,
  		created_at TIMESTAMP DEFAULT now(),
  		updated_at TIMESTAMP DEFAULT now() 
		);
	`
const CREATE_TRANSACTION_STATUS_ENUMS = `
		CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'canceled');
	`

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusConfirmed TransactionStatus = "confirmed"
	StatusCanceled  TransactionStatus = "canceled"
)

type Transaction struct {
	ID                string            `json:"tx_id"`
	FromUserID        string            `json:"from_user_id"`
	ToUserID          string            `json:"to_user_id"`
	Amount            int64             `json:"amount"`
	Status            TransactionStatus `json:"status"`
	ConfirmationToken string            `json:"confirmation_token"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

func NewTransaction(fromUserID, toUserID string, token string, amount int64) *Transaction {
	id := generateTransactionID("NBT")
	return &Transaction{
		ID:                id,
		FromUserID:        fromUserID,
		ToUserID:          toUserID,
		Amount:            amount,
		Status:            StatusPending,
		ConfirmationToken: token,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}

func generateTransactionID(prefix string) string {
	first := 9
	rest := ""
	for range 20 {
		rest += strconv.Itoa(rand.Intn(10))
	}
	return fmt.Sprintf("%s-%d%s", prefix, first, rest)
}
