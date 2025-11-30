package auth

import "errors"

var (
	ErrUserNotFound      = errors.New("auth: user not found")
	ErrEmailAlreadyUsed  = errors.New("auth: email already exists")
	ErrInvalidCredential = errors.New("auth: invalid credential")
	ErrInvalidToken      = errors.New("auth: invalid token")
)