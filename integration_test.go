package constellation_test

import (
	"testing"

	"github.com/tanner-caffrey/constellation-go"
)

// integrationClient is a shared client instance for integration tests
var integrationClient *constellation.Client

func init() {
	// Initialize the integration test client
	integrationClient = constellation.NewClient()
}

// TestGetAPIInfoIntegration tests the GetAPIInfo endpoint with real API
func TestGetAPIInfoIntegration(t *testing.T) {
	info, err := integrationClient.GetAPIInfo()
	if err != nil {
		t.Fatalf("Failed to get API info: %v", err)
	}

	// Basic validation
	if info.DaysIndexed <= 0 {
		t.Errorf("Expected positive days indexed, got %d", info.DaysIndexed)
	}

	t.Logf("API indexed %d days", info.DaysIndexed)
}

// TestGetLinksIntegration tests the GetLinks endpoint with real API
func TestGetLinksIntegration(t *testing.T) {
	// Example target URI for testing
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	linksParams := constellation.LinksParams{
		Target:     target,
		Collection: "app.bsky.feed.like",
		Path:       ".subject.uri",
		Limit:      5,
	}

	links, err := integrationClient.GetLinks(linksParams)
	if err != nil {
		t.Fatalf("Failed to get links: %v", err)
	}

	// Basic validation
	if links.Total < 0 {
		t.Errorf("Expected non-negative total, got %d", links.Total)
	}

	if len(links.LinkingRecords) > linksParams.Limit {
		t.Errorf("Expected at most %d records, got %d", linksParams.Limit, len(links.LinkingRecords))
	}

	t.Logf("Found %d links (total: %d)", len(links.LinkingRecords), links.Total)
}

// TestGetLinksCountIntegration tests the GetLinksCount endpoint with real API
func TestGetLinksCountIntegration(t *testing.T) {
	// Example target URI for testing
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	linksParams := constellation.LinksParams{
		Target:     target,
		Collection: "app.bsky.feed.like",
		Path:       ".subject.uri",
	}

	count, err := integrationClient.GetLinksCount(linksParams)
	if err != nil {
		t.Fatalf("Failed to get links count: %v", err)
	}

	// Basic validation
	if count.Total < 0 {
		t.Errorf("Expected non-negative count, got %d", count.Total)
	}

	t.Logf("Total links count: %d", count.Total)
}

// TestGetDistinctDIDsIntegration tests the GetDistinctDIDs endpoint with real API
func TestGetDistinctDIDsIntegration(t *testing.T) {
	didsParams := constellation.LinksParams{
		Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
		Collection: "app.bsky.graph.block",
		Path:       ".subject",
		Limit:      10,
	}

	dids, err := integrationClient.GetDistinctDIDs(didsParams)
	if err != nil {
		t.Fatalf("Failed to get distinct DIDs: %v", err)
	}

	// Basic validation
	if dids.Total < 0 {
		t.Errorf("Expected non-negative total, got %d", dids.Total)
	}

	if len(dids.DIDs) > didsParams.Limit {
		t.Errorf("Expected at most %d DIDs, got %d", didsParams.Limit, len(dids.DIDs))
	}

	t.Logf("Found %d distinct DIDs (total: %d)", len(dids.DIDs), dids.Total)
}

// TestGetDistinctDIDsCountIntegration tests the GetDistinctDIDsCount endpoint with real API
func TestGetDistinctDIDsCountIntegration(t *testing.T) {
	didsParams := constellation.LinksParams{
		Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
		Collection: "app.bsky.graph.block",
		Path:       ".subject",
	}

	count, err := integrationClient.GetDistinctDIDsCount(didsParams)
	if err != nil {
		t.Fatalf("Failed to get distinct DID count: %v", err)
	}

	// Basic validation
	if count < 0 {
		t.Errorf("Expected non-negative count, got %d", count)
	}

	t.Logf("Found %d distinct DIDs", count)
}

// TestGetLinksWithCursorIntegration tests pagination with cursor
func TestGetLinksWithCursorIntegration(t *testing.T) {
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	// First request
	linksParams := constellation.LinksParams{
		Target:     target,
		Collection: "app.bsky.feed.like",
		Path:       ".subject.uri",
		Limit:      2,
	}

	links1, err := integrationClient.GetLinks(linksParams)
	if err != nil {
		t.Fatalf("Failed to get first page of links: %v", err)
	}

	if len(links1.LinkingRecords) == 0 {
		t.Skip("No links found, skipping pagination test")
	}

	// Second request with cursor
	if links1.Cursor != "" {
		linksParams.Cursor = links1.Cursor
		links2, err := integrationClient.GetLinks(linksParams)
		if err != nil {
			t.Fatalf("Failed to get second page of links: %v", err)
		}

		// Should get different results
		if len(links2.LinkingRecords) > 0 && len(links1.LinkingRecords) > 0 {
			if links1.LinkingRecords[0].RKey == links2.LinkingRecords[0].RKey {
				t.Log("Warning: Cursor pagination may not be working as expected")
			}
		}

		t.Logf("First page: %d records, Second page: %d records",
			len(links1.LinkingRecords), len(links2.LinkingRecords))
	} else {
		t.Log("No cursor returned, pagination not available")
	}
}

// TestGetLinksDifferentCollectionsIntegration tests different collection types
func TestGetLinksDifferentCollectionsIntegration(t *testing.T) {
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	collections := []string{
		"app.bsky.feed.like",
		"app.bsky.feed.repost",
		"app.bsky.graph.follow",
	}

	for _, collection := range collections {
		t.Run(collection, func(t *testing.T) {
			linksParams := constellation.LinksParams{
				Target:     target,
				Collection: collection,
				Path:       ".subject.uri",
				Limit:      3,
			}

			links, err := integrationClient.GetLinks(linksParams)
			if err != nil {
				t.Fatalf("Failed to get links for collection %s: %v", collection, err)
			}

			t.Logf("Collection %s: %d links (total: %d)",
				collection, len(links.LinkingRecords), links.Total)
		})
	}
}
