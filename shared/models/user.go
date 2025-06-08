package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const CREATE_USERS_TABLE = `
	CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT now()
    );
	`

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Passowrd  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username, email, hashedPassowrd string) *User {
	id := generateUserID("NBU")
	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Passowrd:  hashedPassowrd,
		CreatedAt: time.Now().UTC(),
	}
}

func generateUserID(prefix string) string {
	first := 1
	rest := ""
	for range 7 {
		rest += strconv.Itoa(rand.Intn(10))
	}
	return fmt.Sprintf("%s-%d%s", prefix, first, rest)
}
