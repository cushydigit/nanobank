package models

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Passowrd string `json:"-"`
}

func NewUser(username, email, hashedPassowrd string) *User {
	id := fmt.Sprintf("NBU-%s", uuid.New().String())
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Passowrd: hashedPassowrd,
	}
}
