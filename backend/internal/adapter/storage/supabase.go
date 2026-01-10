package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SupabaseStorage struct {
	projectURL string
	serviceKey string
	bucket     string
}

func NewSupabaseStorage() *SupabaseStorage {
	return &SupabaseStorage{
		projectURL: os.Getenv("SUPABASE_PROJECT_URL"),
		serviceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		bucket:     os.Getenv("SUPABASE_BUCKET"),
	}
}

type SignedUploadResponse struct {
	SignedURL string `json:"signedURL"`
	Token     string `json:"token"`
	Path      string `json:"path"`
}

func (s *SupabaseStorage) CreateSignedUploadURL(filename string) (*SignedUploadResponse, error) {
	url := fmt.Sprintf("%s/storage/v1/object/upload/sign/%s/%s",
		s.projectURL, s.bucket, filename)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error (%d): %s", resp.StatusCode, string(body))
	}

	var result SignedUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *SupabaseStorage) GetPublicURL(path string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		s.projectURL, s.bucket, path)
}
