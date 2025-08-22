package constellation

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// LinksParams represents parameters for links-related API calls
type LinksParams struct {
	Target     string // Required: The target URI to find links for
	Collection string // Optional: Filter by collection type
	Path       string // Optional: JSONPath to the target within records
	Limit      int    // Optional: Maximum number of results to return
	Cursor     string // Optional: Cursor for pagination
}

// LinksResponse represents the response from links endpoints
type LinksResponse struct {
	Total          int          `json:"total,omitempty"`
	LinkingRecords []LinkRecord `json:"linking_records,omitempty"`
	Cursor         string       `json:"cursor,omitempty"`
}

// CountResponse represents the response from count endpoints
type CountResponse struct {
	Total int `json:"total"`
}

// DistinctDIDsResponse represents the response from distinct DIDs endpoints
type DistinctDIDsResponse struct {
	Total  int      `json:"total,omitempty"`
	DIDs   []string `json:"linking_dids,omitempty"`
	Cursor string   `json:"cursor,omitempty"`
}

// GetLinks retrieves a list of records linking to a target
// Endpoint: GET /links
func (c *Client) GetLinks(params LinksParams) (*LinksResponse, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	urlParams := url.Values{}
	urlParams.Add("target", params.Target)

	if params.Collection != "" {
		urlParams.Add("collection", params.Collection)
	}
	if params.Path != "" {
		urlParams.Add("path", params.Path)
	}
	if params.Limit > 0 {
		urlParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		urlParams.Add("cursor", params.Cursor)
	}

	resp, err := c.makeRequest("/links", urlParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var linksResp LinksResponse
	if err := json.NewDecoder(resp.Body).Decode(&linksResp); err != nil {
		return nil, fmt.Errorf("failed to decode links response: %w", err)
	}

	return &linksResp, nil
}

// GetLinksCount retrieves the total number of links pointing at a given target
// Endpoint: GET /links/count
func (c *Client) GetLinksCount(params LinksParams) (*CountResponse, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	urlParams := url.Values{}
	urlParams.Add("target", params.Target)

	if params.Collection != "" {
		urlParams.Add("collection", params.Collection)
	}
	if params.Path != "" {
		urlParams.Add("path", params.Path)
	}

	resp, err := c.makeRequest("/links/count", urlParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var countResp CountResponse
	if err := json.NewDecoder(resp.Body).Decode(&countResp); err != nil {
		return nil, fmt.Errorf("failed to decode count response: %w", err)
	}

	return &countResp, nil
}

// GetDistinctDIDs retrieves a list of distinct DIDs linking to a target
// Endpoint: GET /links/distinct-dids
func (c *Client) GetDistinctDIDs(params LinksParams) (*DistinctDIDsResponse, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	urlParams := url.Values{}
	urlParams.Add("target", params.Target)

	if params.Collection != "" {
		urlParams.Add("collection", params.Collection)
	}
	if params.Path != "" {
		urlParams.Add("path", params.Path)
	}
	if params.Limit > 0 {
		urlParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		urlParams.Add("cursor", params.Cursor)
	}

	resp, err := c.makeRequest("/links/distinct-dids", urlParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var didsResp DistinctDIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&didsResp); err != nil {
		return nil, fmt.Errorf("failed to decode distinct DIDs response: %w", err)
	}

	return &didsResp, nil
}

// GetDistinctDIDs retrieves a list of distinct DIDs linking to a target
// Endpoint: GET /links/distinct-dids
func (c *Client) GetDistinctDIDsCount(params LinksParams) (int, error) {
	if params.Target == "" {
		return -1, fmt.Errorf("target parameter is required")
	}

	urlParams := url.Values{}
	urlParams.Add("target", params.Target)

	if params.Collection != "" {
		urlParams.Add("collection", params.Collection)
	}
	if params.Path != "" {
		urlParams.Add("path", params.Path)
	}
	if params.Limit > 0 {
		urlParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		urlParams.Add("cursor", params.Cursor)
	}

	resp, err := c.makeRequest("/links/count/distinct-dids", urlParams)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var didsResp DistinctDIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&didsResp); err != nil {
		return -1, fmt.Errorf("failed to decode distinct DIDs response: %w", err)
	}

	return didsResp.Total, nil
}
