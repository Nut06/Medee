package auth

import "github.com/google/uuid"

// User is the domain representation for authentication flows.
// Keep it free from infrastructure concerns (no JSON or GORM tags).
type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Role         string
	PasswordHash string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


type LoginResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}