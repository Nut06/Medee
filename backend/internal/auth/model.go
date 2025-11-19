package auth

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type RegisterRespone struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Token struct {
	AccessToken string 
	RefreshToken string
}