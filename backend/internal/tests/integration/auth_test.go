package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/domain/auth"
	"backend/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestAuth_RegisterAndLoginFlow(t *testing.T) {
	// Clean the database table before running this integration test
	clearUsersTable()

	// 1. Register User Action
	registerReq := auth.RegisterRequest{
		FirstName: "Integration",
		LastName:  "Test",
		Email:     "integration@example.com",
		Password:  "Password123!",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute request via Fiber's Test method
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the side effect: User should be saved in DB
	var count int64
	database.DB.Table("users").Where("email = ?", "integration@example.com").Count(&count)
	assert.Equal(t, int64(1), count)

	// 2. Login User Action
	loginReq := auth.LoginRequest{
		Email:    "integration@example.com",
		Password: "Password123!",
	}

	loginBody, _ := json.Marshal(loginReq)
	reqLogin := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")

	respLogin, err := app.Test(reqLogin, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLogin.StatusCode)

	// Verify the side effect: App should issue a JWT Token
	var loginResp auth.LoginResponse
	err = json.NewDecoder(respLogin.Body).Decode(&loginResp)
	assert.NoError(t, err)

	assert.NotEmpty(t, loginResp.Token)
	assert.Equal(t, "integration@example.com", loginResp.User.Email)

	reqProfile := httptest.NewRequest(http.MethodGet, "/user", nil)
	reqProfile.Header.Set("Content-Type", "application/json")
	reqProfile.Header.Set("Authorization", "Bearer " +loginResp.Token)

	respProfile, err := app.Test(reqProfile, -1)

	assert.NoError(t, err)
	
	assert.Equal(t, http.StatusOK, respProfile.StatusCode)

	var profileResp auth.UserResponse

	err = json.NewDecoder(respProfile.Body).Decode(&profileResp)
	assert.NoError(t,err)
}