package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/domain/auth"
	userdomain "backend/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helper: register and login a test user, returns the access token
// ─────────────────────────────────────────────────────────────────────────────
func registerAndLogin(t *testing.T, email string) string {
	t.Helper()

	clearUsersTable()

	// Register
	registerBody, _ := json.Marshal(auth.RegisterRequest{
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Password:  "Password123!",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Flush refresh tokens to avoid collision on back-to-back calls
	clearUsersTable()
	registerAndGetToken(t, email)
	return loginGetToken(t, email)
}

func registerAndGetToken(t *testing.T, email string) {
	t.Helper()
	registerBody, _ := json.Marshal(auth.RegisterRequest{
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Password:  "Password123!",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req, -1) //nolint
}

func loginGetToken(t *testing.T, email string) string {
	t.Helper()
	clearRefreshTokens()

	loginBody, _ := json.Marshal(auth.LoginRequest{
		Email:    email,
		Password: "Password123!",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp auth.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)
	return loginResp.Token
}

func authReq(method, path string, body []byte, token string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: UpdateProfile
// ─────────────────────────────────────────────────────────────────────────────
func TestUser_UpdateProfile(t *testing.T) {
	token := registerAndLogin(t, "user.profile@example.com")

	updateBody, _ := json.Marshal(userdomain.UpdateProfileCommand{
		FirstName:   "Updated",
		LastName:    "Name",
		LinkedInURL: "https://linkedin.com/in/test",
	})
	resp, err := app.Test(authReq(http.MethodPut, "/user", updateBody, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the returned data
	var profileResp userdomain.UserProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profileResp)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", profileResp.FirstName)
	assert.Equal(t, "Name", profileResp.LastName)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Add and Delete Experience
// ─────────────────────────────────────────────────────────────────────────────
func TestUser_AddAndDeleteExperience(t *testing.T) {
	token := registerAndLogin(t, "user.exp@example.com")

	startDate := "2022-01-01"
	expBody, _ := json.Marshal(userdomain.AddExperienceCommand{
		Position:    "Software Engineer",
		CompanyName: "Test Corp",
		StartDate:   &startDate,
	})

	// Add experience
	resp, err := app.Test(authReq(http.MethodPost, "/user/experience", expBody, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response to get the experience ID
	var expResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&expResp)
	assert.NoError(t, err)

	expID := fmt.Sprintf("%v", expResp["id"])
	assert.NotEmpty(t, expID)

	// Delete experience
	delResp, err := app.Test(authReq(http.MethodDelete, "/user/experience/"+expID, nil, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, delResp.StatusCode)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Add and Delete Education
// ─────────────────────────────────────────────────────────────────────────────
func TestUser_AddAndDeleteEducation(t *testing.T) {
	token := registerAndLogin(t, "user.edu@example.com")

	eduBody, _ := json.Marshal(userdomain.AddEducationCommand{
		InstituteName:  "Test University",
		Degree:         "Bachelor's",
		FieldOfStudy:   "Computer Science",
		GraduationYear: 2023,
	})

	// Add education
	resp, err := app.Test(authReq(http.MethodPost, "/user/education", eduBody, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response to get ID
	var eduResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&eduResp)
	assert.NoError(t, err)

	eduID := fmt.Sprintf("%v", eduResp["id"])
	assert.NotEmpty(t, eduID)

	// Delete education
	delResp, err := app.Test(authReq(http.MethodDelete, "/user/education/"+eduID, nil, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, delResp.StatusCode)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Get Full Profile includes Experiences
// ─────────────────────────────────────────────────────────────────────────────
func TestUser_GetFullProfile(t *testing.T) {
	token := registerAndLogin(t, "user.fullprofile@example.com")

	// Add one experience first
	startDate := "2021-05-01"
	expBody, _ := json.Marshal(userdomain.AddExperienceCommand{
		Position:    "Backend Developer",
		CompanyName: "Acme Inc",
		StartDate:   &startDate,
	})
	_, _ = app.Test(authReq(http.MethodPost, "/user/experience", expBody, token), -1)

	// Get full profile
	resp, err := app.Test(authReq(http.MethodGet, "/user/profile", nil, token), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fullProfile userdomain.FullProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&fullProfile)
	assert.NoError(t, err)
	assert.Equal(t, "user.fullprofile@example.com", fullProfile.Email)
	assert.GreaterOrEqual(t, len(fullProfile.Experiences), 1)
	assert.Equal(t, "Backend Developer", fullProfile.Experiences[0].Position)
}
