package institute_adapter

import (
	"backend/internal/domain/domain"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type HipoAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHipoAPIClient() *HipoAPIClient {
	baseURL := os.Getenv("HIPO_BASE_URL")
	if baseURL == "" {
		baseURL = "http://universities.hipolabs.com"
	}
	return &HipoAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: utils.FiveSec,
		},
	}
}

type HipoUniversityResponse struct {
	Name          string   `json:"name"`
	Country       string   `json:"country"`
	AlphaTwoCode  string   `json:"alpha_two_code"`
	StateProvince *string  `json:"state-province"`
	Domains       []string `json:"domains"`
	WebPages      []string `json:"web_pages"`
}

func (c *HipoAPIClient) SearchHipoAPI(ctx context.Context, query string, country string) ([]domain.Institute, error) {
	params := url.Values{}
	params.Add("name", query)
	if country != "" {
		params.Add("country", country)
	}

	apiURL := fmt.Sprintf("%s/search?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hipo API returned status: %d", resp.StatusCode)
	}

	var hipoResults []HipoUniversityResponse
	if err := json.NewDecoder(resp.Body).Decode(&hipoResults); err != nil {
		return nil, err
	}

	// Map to domain entities
	institutes := make([]domain.Institute, len(hipoResults))
	for i, hr := range hipoResults {
		var domainStr *string
		if len(hr.Domains) > 0 {
			domainStr = &hr.Domains[0]
		}

		institutes[i] = domain.Institute{
			Name:     hr.Name,
			Country:  hr.Country,
			Domain:   domainStr,
			WebPages: hr.WebPages,
			Source:   "hipo_api",
		}
	}

	return institutes, nil
}
