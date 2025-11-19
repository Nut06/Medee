package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRespone struct {
	ID    uuid.UUID `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID uuid.UUID `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
}