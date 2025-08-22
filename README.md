# Constellation API Client

A Go client library for interfacing with the [Constellation API](https://constellation.microcosm.blue/) - a service that provides access to Bluesky's public data via a simple REST API.

## Features

- **Complete API Coverage**: Implements all available Constellation API endpoints
- **Type-Safe**: Strongly typed request/response structures
- **Error Handling**: Comprehensive error handling with descriptive messages
- **Configurable**: Customizable base URL and HTTP client settings

## Installation

```bash
go get github.com/tanner-caffrey/constellation-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/tanner-caffrey/constellation-go"
)

func main() {
    // Create a new client
    client := constellation.NewClient()
    
    // Get API information
    info, err := client.GetAPIInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Days indexed: %d\n", info.DaysIndexed)
    
    // Get links to a target
    params := constellation.LinksParams{
        Target:     "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r",
        Collection: "app.bsky.feed.like",
        Path:       ".subject.uri",
        Limit:      5,
    }
    
    links, err := client.GetLinks(params)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d links\n", len(links.LinkingRecords))
}
```

## API Reference

### Client Creation

```go
// Create client with default settings
client := constellation.NewClient()

// Create client with custom configuration
client := constellation.NewClientWithConfig(
    "https://custom-api-url.com",
    60*time.Second, // timeout
)

// Create client with custom User-Agent
client := constellation.NewClientWithUserAgent("my-app/1.0.0")
```

### User-Agent Configuration

The client supports multiple ways to configure the User-Agent string:

#### 1. Environment Variable (Recommended)
Set the `CONSTELLATION_USER_AGENT` environment variable:
```bash
export CONSTELLATION_USER_AGENT="my-app-constellation-client/1.0.0"
```

Or in a `.env` file:
```env
CONSTELLATION_USER_AGENT=my-app-constellation-client/1.0.0
```

#### 2. Custom User-Agent Method
Use `NewClientWithUserAgent()` to set a custom User-Agent:
```go
client := constellation.NewClientWithUserAgent("my-custom-user-agent/2.0.0")
```

#### 3. Default User-Agent
If no environment variable is set, the default User-Agent is `constellation-go/1.0.0`.

**Priority Order:**
1. `NewClientWithUserAgent()` (highest priority)
2. `CONSTELLATION_USER_AGENT` environment variable
3. Default User-Agent (lowest priority)

### Available Methods

#### GetAPIInfo()
Get basic information about the Constellation API including statistics.

```go
info, err := client.GetAPIInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("API indexed %d days\n", info.DaysIndexed)
```

#### GetLinks(params LinksParams)
Retrieve records that link to a specific target.

```go
params := constellation.LinksParams{
    Target:     "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r",
    Collection: "app.bsky.feed.like", // optional
    Path:       ".subject.uri",       // optional
    Limit:      5,                    // optional
    Cursor:     "",                   // optional for pagination
}
links, err := client.GetLinks(params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d links (total: %d)\n", len(links.LinkingRecords), links.Total)
```

#### GetLinksCount(params LinksParams)
Get the total count of links pointing to a target.

```go
params := constellation.LinksParams{
    Target:     "at://did:plc:vc7f4oafdgxsihk4cry2xpze/app.bsky.feed.post/3lgwdn7vd722r",
    Collection: "app.bsky.feed.like",
    Path:       ".subject.uri",
}
count, err := client.GetLinksCount(params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total links: %d\n", count.Total)
```

#### GetDistinctDIDs(params LinksParams)
Get a list of unique DIDs that link to a target.

```go
params := constellation.LinksParams{
    Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
    Collection: "app.bsky.graph.block",
    Path:       ".subject",
    Limit:      10,
}
dids, err := client.GetDistinctDIDs(params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Unique DIDs: %d\n", len(dids.DIDs))
```

#### GetDistinctDIDsCount(params LinksParams)
Get the total count of distinct DIDs linking to a target.

```go
params := constellation.LinksParams{
    Target:     "did:plc:vc7f4oafdgxsihk4cry2xpze",
    Collection: "app.bsky.graph.block",
    Path:       ".subject",
}
count, err := client.GetDistinctDIDsCount(params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total distinct DIDs: %d\n", count)
```

## Data Structures

### LinksParams
Parameters for links-related API calls:
- `Target` (required): The target URI to find links for
- `Collection` (optional): Filter by collection type
- `Path` (optional): JSONPath to the target within records
- `Limit` (optional): Maximum number of results
- `Cursor` (optional): Pagination cursor

### LinkRecord
Represents a link record from the API:
- `DID`: The DID of the record author
- `Collection`: The collection type
- `RKey`: The record key
- `URI`: The full AT-URI
- `CID`: The content identifier
- `IndexedAt`: When the record was indexed
- `Value`: The record content

### APIResponse
Response from the GetAPIInfo endpoint:
- `DaysIndexed`: Number of days the API has been indexing data
- `Stats`: Statistics about the indexed data
- `Help`: Help information (if available)

### LinksResponse
Response from GetLinks endpoint:
- `Total`: Total number of matching records
- `LinkingRecords`: Array of link records
- `Cursor`: Pagination cursor for next page

### CountResponse
Response from count endpoints:
- `Total`: Total count of matching records

### DistinctDIDsResponse
Response from GetDistinctDIDs endpoint:
- `Total`: Total number of distinct DIDs
- `DIDs`: Array of distinct DID strings
- `Cursor`: Pagination cursor for next page

## Error Handling

All methods return an error as the second return value. Common error scenarios include:

- Network connectivity issues
- Invalid parameters (e.g., missing required `Target`)
- API rate limiting
- Invalid API responses

```go
links, err := client.GetLinks(params)
if err != nil {
    log.Printf("Error fetching links: %v", err)
    return
}
```

## Contributing

This package is designed to be a complete interface to the Constellation API. If you notice missing functionality or bugs, please open an issue or submit a pull request.

## Acknowledgments

This Go client library is an unofficial interface to the [Constellation API](https://constellation.microcosm.blue/), which is developed and maintained by the [Microcosm](https://microcosm.blue/) team. 

**Important Note**: This is an independent Go client library and is not officially affiliated with or endorsed by Microcosm or the Constellation API team. It is simply a Go interface to their excellent API service I created for my own projects.
