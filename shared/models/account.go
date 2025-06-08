package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const CREATE_ACCOUNTS_TABLE = `
		CREATE TABLE IF NOT EXISTS accounts (
  		id TEXT PRIMARY KEY,
  		user_id TEXT UNIQUE NOT NULL,
  		username TEXT NOT NULL,
  		balance BIGINT NOT NULL DEFAULT 0,
  		created_at TIMESTAMP DEFAULT now(),
  		updated_at TIMESTAMP DEFAULT now()
		);
	`

type Account struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Balance   int64     `json:"balance"` // fixed scale (cents)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAccount(userID, username string) *Account {
	id := generateAccountNumber("NBA")
	return &Account{
		ID:        id,
		UserID:    userID,
		Username:  username,
		Balance:   0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func generateAccountNumber(prefix string) string {
	first := 9910
	rest := ""
	for range 12 {
		rest += strconv.Itoa(rand.Intn(10))
	}
	return fmt.Sprintf("%s-%d%s", prefix, first, rest)
}
