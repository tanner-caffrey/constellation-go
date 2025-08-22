package constellation_test

import (
	"testing"

	"github.com/tanner-caffrey/constellation-go"
)

// testClient is a shared client instance for tests
var testClient *constellation.Client

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Initialize the test client
	testClient = constellation.NewClientWithUserAgent("constellation-go-test/1.0.0")

	// Run the tests
	m.Run()
}

// TestGetAPIInfo tests the GetAPIInfo endpoint
func TestGetAPIInfo(t *testing.T) {
	info, err := testClient.GetAPIInfo()
	if err != nil {
		t.Fatalf("Failed to get API info: %v", err)
	}

	// Basic validation
	if info.DaysIndexed <= 0 {
		t.Errorf("Expected positive days indexed, got %d", info.DaysIndexed)
	}

	t.Logf("API indexed %d days", info.DaysIndexed)
}

// TestGetLinks tests the GetLinks endpoint
func TestGetLinks(t *testing.T) {
	// Example target URI for testing
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	linksParams := constellation.LinksParams{
		Target:     target,
		Collection: "app.bsky.feed.like",
		Path:       ".subject.uri",
		Limit:      5,
	}

	links, err := testClient.GetLinks(linksParams)
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

// TestGetLinksCount tests the GetLinksCount endpoint
func TestGetLinksCount(t *testing.T) {
	// Example target URI for testing
	target := "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r"

	linksParams := constellation.LinksParams{
		Target:     target,
		Collection: "app.bsky.feed.like",
		Path:       ".subject.uri",
	}

	count, err := testClient.GetLinksCount(linksParams)
	if err != nil {
		t.Fatalf("Failed to get links count: %v", err)
	}

	// Basic validation
	if count.Total < 0 {
		t.Errorf("Expected non-negative count, got %d", count.Total)
	}

	t.Logf("Total links count: %d", count.Total)
}

// TestGetDistinctDIDs tests the GetDistinctDIDs endpoint
func TestGetDistinctDIDs(t *testing.T) {
	didsParams := constellation.LinksParams{
		Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
		Collection: "app.bsky.graph.block",
		Path:       ".subject",
		Limit:      10,
	}

	dids, err := testClient.GetDistinctDIDs(didsParams)
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

// TestGetDistinctDIDsCount tests the GetDistinctDIDsCount endpoint
func TestGetDistinctDIDsCount(t *testing.T) {
	didsParams := constellation.LinksParams{
		Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
		Collection: "app.bsky.graph.block",
		Path:       ".subject",
	}

	count, err := testClient.GetDistinctDIDsCount(didsParams)
	if err != nil {
		t.Fatalf("Failed to get distinct DID count: %v", err)
	}

	// Basic validation
	if count < 0 {
		t.Errorf("Expected non-negative count, got %d", count)
	}

	t.Logf("Found %d distinct DIDs", count)
}
