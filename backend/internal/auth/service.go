package auth

import (
	"backend/internal/database"
	"backend/internal/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func register(input *RegisterRequest) (*RegisterRespone, error) {
	var existing *user.User
	existing, err := FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("Email is already exist")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := string(user.RoleCandidate)
	hpwd := string(hashed)
	newUser := &user.User{
		Name:     input.Name,
		Email:    input.Email,
		Role:     &role,
		Password: &hpwd,
	}

	if err := user.CreateUser(newUser); err != nil {
		return nil, err
	}

	createdUser, _ := FindByEmail(input.Email);

	return &RegisterRespone{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  string(*createdUser.Role),
	}, nil
}


func login(input *LoginRequest) (*LoginResponse, *Token, error) {
	foundUser, err := FindByEmail(input.Email)
	if err != nil {
		return nil, nil, errors.New("No existing user")
	}
	userPassword := *foundUser.Password
	if err := comparePassword(userPassword, input.Password); err != nil {
		return nil, nil, errors.New("Invalid Username or Password")
	}

	accessToken, err := GenerateAccessToken(foundUser.ID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken := GenerateRefreshToken()
	return &LoginResponse{
			UserID: foundUser.ID,
			Name:   foundUser.Name,
			Email:  foundUser.Email,
			Role:   string(*foundUser.Role),
		}, &Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, err
}

func comparePassword(hpw string, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpw), []byte(pw))
}

func FindByEmail(email string) (*user.User, error) {

	var user user.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func logout(token string) error {
	return database.DB.Delete("token = ?", token).Delete(&user.RefreshToken{}).Error
}
