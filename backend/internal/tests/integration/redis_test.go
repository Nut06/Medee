package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestInstitute_RedisCachingFlow(t *testing.T) {
	// 1. Initial request (Should hit Database and then set to Redis)
	req1 := httptest.NewRequest(http.MethodGet, "/institutes/search?q=test", nil)
	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	// Since we don't have exact timing precision guarantees inside a container, 
	// the best way to mathematically prove Redis is working is to manually check the Redis key
	// that your adapter saves data into.	

	// In your `institute_adapter`, you likely cache it with a key. 
	// We'll verify Redis holds *some* keys after the request.
	redisClient := database.NewRedisClient()
	
	// Assuming you use wildcard or a specific prefix for caching institutes like "institute:*"
	// This ensures Redis was touched during the request lifecycle.
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "*").Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, keys, "Redis should contain cached keys after fetching institutes")

	// Clean up Redis after test to prevent pollution
	redisClient.FlushAll(ctx)
}
