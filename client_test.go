package constellation_test

import (
	"os"
	"testing"
	"time"

	"github.com/tanner-caffrey/constellation-go"
)

// TestClientCreation tests client creation with different configurations
func TestClientCreation(t *testing.T) {
	// Test default client
	client1 := constellation.NewClient()
	if client1 == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client1.BaseURL != constellation.DefaultBaseURL {
		t.Errorf("Expected BaseURL '%s', got '%s'", constellation.DefaultBaseURL, client1.BaseURL)
	}

	// Test client with custom config
	customURL := "https://example.com"
	customTimeout := 60 * time.Second
	client2 := constellation.NewClientWithConfig(customURL, customTimeout)
	if client2 == nil {
		t.Fatal("NewClientWithConfig() returned nil")
	}

	if client2.BaseURL != customURL {
		t.Errorf("Expected BaseURL '%s', got '%s'", customURL, client2.BaseURL)
	}
}

// TestUserAgentEnvironmentVariable tests User-Agent from environment variable
func TestUserAgentEnvironmentVariable(t *testing.T) {
	// Test with environment variable set
	customUserAgent := "my-custom-app/2.0.0"
	os.Setenv(constellation.EnvUserAgent, customUserAgent)
	defer os.Unsetenv(constellation.EnvUserAgent)

	client := constellation.NewClient()
	if client.UserAgent != customUserAgent {
		t.Errorf("Expected User-Agent '%s', got '%s'", customUserAgent, client.UserAgent)
	}

	// Test with environment variable unset
	os.Unsetenv(constellation.EnvUserAgent)
	client2 := constellation.NewClient()
	if client2.UserAgent != constellation.DefaultUserAgent {
		t.Errorf("Expected User-Agent '%s', got '%s'", constellation.DefaultUserAgent, client2.UserAgent)
	}
}

// TestUserAgentOverride tests that NewClientWithUserAgent overrides environment variable
func TestUserAgentOverride(t *testing.T) {
	// Set environment variable
	customUserAgent := "env-user-agent/1.0.0"
	os.Setenv(constellation.EnvUserAgent, customUserAgent)
	defer os.Unsetenv(constellation.EnvUserAgent)

	// Test that NewClientWithUserAgent overrides environment variable
	overrideUserAgent := "override-user-agent/3.0.0"
	client := constellation.NewClientWithUserAgent(overrideUserAgent)
	if client.UserAgent != overrideUserAgent {
		t.Errorf("Expected User-Agent '%s', got '%s'", overrideUserAgent, client.UserAgent)
	}
}

// TestLinksParamsValidation tests parameter validation logic
func TestLinksParamsValidation(t *testing.T) {
	// Test with empty target (should fail)
	params := constellation.LinksParams{
		Target:     "",
		Collection: "app.bsky.feed.like",
	}

	// Create a client that won't make actual HTTP requests
	client := constellation.NewClientWithConfig("http://invalid-url", 1*time.Second)

	_, err := client.GetLinks(params)
	if err == nil {
		t.Error("Expected error for empty target, got nil")
	} else {
		// Should fail with validation error, not network error
		if err.Error() != "target parameter is required" {
			t.Errorf("Expected validation error, got: %v", err)
		}
	}

	// Test with valid target (should pass validation)
	params.Target = "at://did:plc:example/app.bsky.feed.post/example"
	_, err = client.GetLinks(params)
	// This will fail with network error, but that's expected since we're using an invalid URL
	if err != nil && err.Error() == "target parameter is required" {
		t.Error("Expected network error, got validation error")
	}
}

// TestStructDefinitions tests that the struct definitions are correct
func TestStructDefinitions(t *testing.T) {
	// Test LinksResponse struct
	linksResp := constellation.LinksResponse{
		Total:          100,
		LinkingRecords: []constellation.LinkRecord{},
		Cursor:         "test-cursor",
	}

	if linksResp.Total != 100 {
		t.Errorf("Expected Total 100, got %d", linksResp.Total)
	}

	if linksResp.Cursor != "test-cursor" {
		t.Errorf("Expected Cursor 'test-cursor', got '%s'", linksResp.Cursor)
	}

	// Test LinkRecord struct
	linkRecord := constellation.LinkRecord{
		DID:        "did:plc:example",
		Collection: "app.bsky.feed.like",
		RKey:       "example-rkey",
		URI:        "at://did:plc:example/app.bsky.feed.like/example-rkey",
		CID:        "example-cid",
		IndexedAt:  "2023-01-01T00:00:00Z",
		Value:      map[string]any{"key": "value"},
	}

	if linkRecord.DID != "did:plc:example" {
		t.Errorf("Expected DID 'did:plc:example', got '%s'", linkRecord.DID)
	}

	if linkRecord.Collection != "app.bsky.feed.like" {
		t.Errorf("Expected Collection 'app.bsky.feed.like', got '%s'", linkRecord.Collection)
	}
}
