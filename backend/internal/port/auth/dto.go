package authport

import "backend/internal/domain/auth"

type RegisterCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type RegisterResult struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	TokenPair *auth.TokenPair
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Companies []auth.Company
}

type RefreshCommand struct {
	RefreshToken string
}
