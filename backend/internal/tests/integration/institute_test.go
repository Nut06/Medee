package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test: Institute Search - Happy Path
// ─────────────────────────────────────────────────────────────────────────────
func TestInstitute_SearchWithValidQuery(t *testing.T) {
	// This is a public endpoint, no auth required
	req := httptest.NewRequest(http.MethodGet, "/institutes/search?q=ma", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	// Response shape should always contain "results" and "count"
	_, hasResults := body["results"]
	_, hasCount := body["count"]
	assert.True(t, hasResults, "Response should have a 'results' key")
	assert.True(t, hasCount, "Response should have a 'count' key")
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Institute Search - Query too short (< 2 chars)
// ─────────────────────────────────────────────────────────────────────────────
func TestInstitute_SearchWithShortQueryReturns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/institutes/search?q=a", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Institute Search - Missing query parameter
// ─────────────────────────────────────────────────────────────────────────────
func TestInstitute_SearchWithNoQueryReturns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/institutes/search", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Get Institute by ID - Not Found
// ─────────────────────────────────────────────────────────────────────────────
func TestInstitute_GetByNonExistentIDReturns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/institutes/nonexistent-id-999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// ─────────────────────────────────────────────────────────────────────────────
// Test: Field of Study Search - Happy Path
// ─────────────────────────────────────────────────────────────────────────────
func TestFieldOfStudy_SearchWithValidQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/field-of-studies/search?q=cs", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Should return 200 with results (even if empty, it should not fail)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
}
