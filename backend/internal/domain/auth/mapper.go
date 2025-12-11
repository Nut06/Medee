package auth

import "backend/internal/domain/domain"

// ToUserResponse converts domain.User to auth.UserResponse
func ToUserResponse(u *domain.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:          u.ID.String(),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		PhoneNumber: u.PhoneNumber,
		LinkedInURL: u.LinkedInURL,
		GitHubURL:   u.GitHubURL,
		WebsiteURL:  u.WebsiteURL,
	}
}
