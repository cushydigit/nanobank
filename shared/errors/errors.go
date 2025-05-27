package errors

import "errors"

var (
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
