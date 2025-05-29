package errors

import "errors"

var (
	ErrDuplicateEmail                = errors.New("email already exists")
	ErrUserNotFound                  = errors.New("user not found")
	ErrInternalServer                = errors.New("internal server error")
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrInvalidRefreshToken           = errors.New("invalid refresh token")
	ErrAccountNotFound               = errors.New("account not found")
	ErrAccountAlreadyExists          = errors.New("account already exists")
	ErrAmountMustBePositive          = errors.New("amount must be positive")
	ErrInsufficientBalance           = errors.New("insufficient funds")
	ErrMissingAuthorizationHeader    = errors.New("authorizaion header missing")
	ErrInvalidTokenFormat            = errors.New("invalid token format")
	ErrInvalidOrExpiredToken         = errors.New("invalid or expired token")
	ErrPermissionDenied              = errors.New("permission denied")
	ErrContextValueNotFoundInRequest = errors.New("object not found in context of request")
)
