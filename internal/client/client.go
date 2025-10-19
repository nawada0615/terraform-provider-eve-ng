// Package client provides HTTP client functionality for EVE-NG API interactions.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents the EVE-NG API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	session    string
	username   string
	password   string
}

// Config holds the client configuration
type Config struct {
	Endpoint           string
	Username           string
	Password           string
	InsecureSkipVerify bool
	Timeout            time.Duration
}

// NewClient creates a new EVE-NG API client
func NewClient(config *Config) (*Client, error) {
	baseURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// Ensure base URL ends with /
	if baseURL.Path == "" {
		baseURL.Path = "/"
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// #nosec G402 -- InsecureSkipVerify is configurable for development/test environments
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		},
	}

	client := &Client{
		baseURL:    baseURL.String(),
		httpClient: httpClient,
		username:   config.Username,
		password:   config.Password,
	}

	// Authenticate on creation
	if err := client.Login(); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

// Login authenticates with the EVE-NG API
func (c *Client) Login() error {
	loginData := map[string]interface{}{
		"username": c.username,
		"password": c.password,
		"html5":    1,
	}

	reqBody, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"api/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform login request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for better error messages
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Code    int    `json:"code"`
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("login failed: %s (code: %d)", errorResp.Message, errorResp.Code)
		}
		return fmt.Errorf("login failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Extract session cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "unetlab_session" {
			c.session = cookie.Value
			break
		}
	}

	if c.session == "" {
		return fmt.Errorf("no session cookie received")
	}

	return nil
}

// Logout logs out from the EVE-NG API
func (c *Client) Logout() error {
	req, err := http.NewRequest("GET", c.baseURL+"api/auth/logout", http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:  "unetlab_session",
		Value: c.session,
		Path:  "/api/",
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform logout request: %w", err)
	}
	defer resp.Body.Close()

	c.session = ""
	return nil
}

// Do performs an HTTP request with session authentication
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Add session cookie to all requests
	req.AddCookie(&http.Cookie{
		Name:  "unetlab_session",
		Value: c.session,
		Path:  "/api/",
	})

	return c.httpClient.Do(req)
}

// Get performs a GET request
func (c *Client) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, http.NoBody)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest("PUT", c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+path, http.NoBody)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// HandleResponse handles API responses and extracts data
func (c *Client) HandleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Code    int    `json:"code"`
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("API error %d: %s", errorResp.Code, errorResp.Message)
		}
		return fmt.Errorf("API error with status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
