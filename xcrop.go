// Package xcrop provides a Go client for the XCROP API — an X/Twitter
// data intelligence platform offering user profiles, tweets, search,
// trending topics, and write operations.
//
// Usage:
//
//	client := xcrop.NewClient("xc_live_...")
//	user, err := client.Users.Get(ctx, "elonmusk")
//
// The client automatically retries on rate-limit (429) and server errors (5xx)
// with exponential backoff.
package xcrop

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default XCROP API base URL.
	DefaultBaseURL = "https://xcrop.io/api/v2"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// Version is the SDK version.
	Version = "1.1.0"
)

// Client is the XCROP API client. Use NewClient to create one.
type Client struct {
	// Services
	Users       *UsersService
	Tweets      *TweetsService
	Search      *SearchService
	Lists       *ListsService
	Communities *CommunitiesService
	Trending    *TrendingService
	Account     *AccountService
	Stream      *StreamService

	http *httpClient
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL (e.g., for testing or self-hosted instances).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.http.baseURL = url
	}
}

// WithHTTPClient sets a custom *http.Client for all API requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.http.client = hc
	}
}

// WithTimeout sets the HTTP client timeout. Default is 30 seconds.
// This stores the timeout separately and creates a new internal HTTP client,
// avoiding mutation of any user-provided *http.Client. (#2)
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http.timeout = d
	}
}

// WithMaxRetries sets the maximum number of retries on 429/5xx errors.
// Default is 3. Set to 0 to disable retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		c.http.maxRetries = n
	}
}

// NewClient creates a new XCROP API client with the given API key.
//
//	client := xcrop.NewClient("xc_live_a2bed1ff...")
//	client := xcrop.NewClient("xc_live_...", xcrop.WithTimeout(10*time.Second))
func NewClient(apiKey string, opts ...Option) *Client {
	hc := &httpClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		maxRetries: defaultMaxRetries,
		timeout:    DefaultTimeout,
	}

	c := &Client{http: hc}

	for _, opt := range opts {
		opt(c)
	}

	// Apply stored timeout to the internal client (not mutating user's client) (#2)
	if hc.timeout != DefaultTimeout || hc.client.Timeout != hc.timeout {
		hc.client = &http.Client{
			Timeout:   hc.timeout,
			Transport: hc.client.Transport,
		}
	}

	// Initialize services
	c.Users = newUsersService(hc)
	c.Tweets = newTweetsService(hc)
	c.Search = newSearchService(hc)
	c.Lists = newListsService(hc)
	c.Communities = newCommunitiesService(hc)
	c.Trending = newTrendingService(hc)
	c.Account = newAccountService(hc)
	c.Stream = newStreamService(hc)

	return c
}
