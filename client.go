// Package constellation provides a client for interfacing with the Constellation API
// at https://constellation.microcosm.blue/
package constellation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Constellation API
	DefaultBaseURL = "https://constellation.microcosm.blue"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
	// DefaultUserAgent is the default User-Agent string for API requests
	DefaultUserAgent = "constellation-go/1.0.0"
	// EnvUserAgent is the environment variable name for custom User-Agent
	EnvUserAgent = "CONSTELLATION_USER_AGENT"
)

// getUserAgent returns the User-Agent string, checking environment variable first
func getUserAgent() string {
	if envUserAgent := os.Getenv(EnvUserAgent); envUserAgent != "" {
		return envUserAgent
	}
	return DefaultUserAgent
}

// Client represents a Constellation API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// NewClient creates a new Constellation API client with default settings
func NewClient() *Client {
	return &Client{
		BaseURL:   DefaultBaseURL,
		UserAgent: getUserAgent(),
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithConfig creates a new Constellation API client with custom configuration
func NewClientWithConfig(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL:   baseURL,
		UserAgent: getUserAgent(),
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewClientWithUserAgent creates a new client with a custom User-Agent
func NewClientWithUserAgent(userAgent string) *Client {
	return &Client{
		BaseURL:   DefaultBaseURL,
		UserAgent: userAgent,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// APIResponse represents a generic API response structure
type APIResponse struct {
	Help        string `json:"help,omitempty"`
	DaysIndexed int    `json:"days_indexed,omitempty"`
	Stats       Stats  `json:"stats,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Stats represents the statistics from the API
type Stats struct {
	DIDs           int64 `json:"dids"`
	Targetables    int64 `json:"targetables"`
	LinkingRecords int64 `json:"linking_records"`
}

// LinkRecord represents a single link record from the API
type LinkRecord struct {
	DID        string         `json:"did"`
	Collection string         `json:"collection"`
	RKey       string         `json:"rkey"`
	URI        string         `json:"uri"`
	CID        string         `json:"cid"`
	IndexedAt  string         `json:"indexedAt"`
	Value      map[string]any `json:"value"`
}

// makeRequest performs an HTTP GET request to the specified endpoint with parameters
func (c *Client) makeRequest(endpoint string, params url.Values) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	if len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return resp, nil
}

// GetAPIInfo retrieves basic information about the Constellation API
func (c *Client) GetAPIInfo() (*APIResponse, error) {
	resp, err := c.makeRequest("/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp, nil
}
