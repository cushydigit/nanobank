package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAccount(UserID string) *Account {
	id := generateAccountNumber("NBA")
	return &Account{
		ID:        id,
		UserID:    UserID,
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
