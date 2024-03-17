// client.go

package client

import (
	"fmt"
	"net/http"
)

// Client represents an API client
type Client struct {
	BaseURL string
}

// NewClient creates a new instance of the API client
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

// Get performs a GET request to the API
func (c *Client) Get(endpoint string, authToken string) (*http.Response, error) {
	var client = http.DefaultClient
	client.Transport = &http.Transport{}

	url := c.BaseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	fmt.Println(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	response, err := client.Do(req)

	defer response.Body.Close()

	return response, nil
}
