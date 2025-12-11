package auth

import "github.com/google/uuid"

// User is the domain representation for authentication flows.
// Keep it free from infrastructure concerns (no JSON or GORM tags).
type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Role         string
	PasswordHash string
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type RegisterResponse struct {
	User      UserResponse `json:"user"`
	Companies []Company    `json:"companies"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User      UserResponse `json:"user"`
	Companies []Company    `json:"companies"`
}

type UserResponse struct {
	// Basic Info
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`

	// Profile Info
	AvatarURL   *string `json:"AvatarURL,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`

	// Social Links
	LinkedInURL *string `json:"linkedInURL,omitempty"`
	GitHubURL   *string `json:"githubURL,omitempty"`
	WebsiteURL  *string `json:"websiteURL,omitempty"`
}
