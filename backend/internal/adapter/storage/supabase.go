package storage

import (
	"fmt"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorage struct {
	client *storage_go.Client
	bucket string
}

func NewSupabaseStorage() *SupabaseStorage {
	projectURL := os.Getenv("SUPABASE_PROJECT_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	bucket := os.Getenv("SUPABASE_BUCKET")

	// Initialize storage-go client
	// storage-go expects the storage URL (e.g., scheme://<project_ref>.supabase.co/storage/v1)
	storageURL := fmt.Sprintf("%s/storage/v1", projectURL)
	client := storage_go.NewClient(storageURL, serviceKey, nil)

	return &SupabaseStorage{
		client: client,
		bucket: bucket,
	}
}

type SignedUploadResponse struct {
	SignedURL string `json:"signedURL"`
	Path  string `json:"path"`
}

func (s *SupabaseStorage) CreateSignedUploadURL(filename string) (*SignedUploadResponse, error) {
	resp, err := s.client.CreateSignedUploadUrl(s.bucket, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed upload url: %w", err)
	}
	
	fmt.Print("response upload url success")
	return &SignedUploadResponse{
		SignedURL: resp.Url,
		Path: filename,
	}, nil
}

func (s *SupabaseStorage) GetPublicURL(path string) string {
	return s.client.GetPublicUrl(s.bucket, path).SignedURL
}
