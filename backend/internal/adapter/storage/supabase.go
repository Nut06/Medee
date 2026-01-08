package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SupabaseStorage struct {
	projectURL string
	serviceKey string
	bucket string
}

func NewSupabaseStorage() *SupabaseStorage {
	return &SupabaseStorage{
		projectURL: os.Getenv("SUPABASE_PROJECT_URL"),
		serviceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		bucket: os.Getenv("SUPABASE_BUCKET"),
	}
}

type SingnedUploadResponse struct {
	SignedURL string `json:"signed_url"`
	Token string `json:"token"`
	Path string `json:"path"`
}

func (s *SupabaseStorage) CreateSignedUploadURL(filename string) (*SingnedUploadResponse, error) {
		url := fmt.Sprintf("%s/storage/v1/object/upload/sign/%s/%s",
        s.projectURL, s.bucket, filename)

		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.serviceKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var result SingnedUploadResponse
		json.NewDecoder(resp.Body).Decode(&result)

		return &result, nil
}


func (s *SupabaseStorage) GetPublicURL(path string) string {
	return  fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
        s.projectURL, s.bucket, path)
}