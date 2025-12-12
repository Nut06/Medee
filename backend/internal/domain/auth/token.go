package auth

import "time"

type TokenPair struct {
	AccessToken     string
	RefreshToken    string
	RefreshExpiresAt time.Time
}