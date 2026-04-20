package auth

import "errors"

var (
	ErrUserNotFound        = errors.New("auth: user not found")
	ErrEmailAlreadyUsed    = errors.New("auth: email already exists")
	ErrInvalidCredential   = errors.New("auth: invalid credential")
	ErrInvalidToken        = errors.New("auth: invalid token")
	ErrInvalidRefreshToken = errors.New("auth: invalid refresh token")
	ErrInvalidAccessToken  = errors.New("auth: invalid access token")
	ErrInvalidUser         = errors.New("auth: invalid user")
	ErrInvalidRequestBody  = errors.New("auth: invalid request body")
)

var (
	ErrInvalidIdempotencyKey = errors.New("Wrong key idempotency")
)