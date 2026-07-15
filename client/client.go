package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is the client for Meteora's DLMM REST API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option is a function to configure the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient returns a new Client.
func NewClient(baseURL string, opts ...Option) *Client {
	if baseURL == "" {
		baseURL = "https://dlmm.datapi.meteora.ag"
	}
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doRequest performs an HTTP request and decodes the response into the target interface.
func (c *Client) doRequest(ctx context.Context, path string, query url.Values, target interface{}) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Message != "" {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("API returned status: %d %s", resp.StatusCode, resp.Status)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
