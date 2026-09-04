package api

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"time"
)


type Client struct {
	HTTPClient *http.Client
	APIKey     string
}


func NewClient(apiKey string) *Client {
	return &Client{
		HTTPClient: &http.Client{ Timeout: 1*time.Second },
		APIKey: apiKey,
	}
}

func (c *Client) APICall(path string, query map[string]string) ([]byte, error) {
	baseURL := url.URL{
		Scheme: "https",
		Host: "api.meteo-concept.com",
		Path: "/api/"+path,
	}

	q := baseURL.Query()
	for key, value := range query {
		q.Set(key, value)
	}

	baseURL.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network call error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error with status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return nil, fmt.Errorf("API response reading cancelled: %w", err)
	}
	
	return body, nil
}