package xcrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultMaxRetries  = 3
	defaultRetryBaseMs = 500
	defaultRetryMaxMs  = 30000
)

// decodeError is a non-retryable error from JSON decoding.
type decodeError struct {
	err error
}

func (e *decodeError) Error() string { return e.err.Error() }
func (e *decodeError) Unwrap() error { return e.err }

// httpClient wraps HTTP request execution with retry logic.
type httpClient struct {
	client     *http.Client
	baseURL    string
	apiKey     string
	maxRetries int
	timeout    time.Duration
}

// requestOptions configures a single HTTP request.
type requestOptions struct {
	method string
	path   string
	query  map[string]string
	body   interface{}
}

// do executes an HTTP request with automatic retry on 429 and 5xx errors.
// It decodes the JSON response into dest.
func (h *httpClient) do(ctx context.Context, opts requestOptions, dest interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := h.calcBackoff(attempt, lastErr)
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}

		err := h.doOnce(ctx, opts, dest)
		if err == nil {
			return nil
		}

		// Don't retry JSON decode errors — the response was received successfully
		// but couldn't be parsed, retrying won't help.
		if _, isDecodeErr := err.(*decodeError); isDecodeErr {
			return err.(*decodeError).err
		}

		apiErr, isAPIErr := err.(*APIError)
		if !isAPIErr {
			// Network error — retry
			lastErr = err
			continue
		}

		// Only retry on 429 or 5xx
		if apiErr.StatusCode == 429 || (apiErr.StatusCode >= 500 && apiErr.StatusCode < 600) {
			lastErr = err
			continue
		}

		// Non-retryable API error
		return err
	}

	return lastErr
}

// doOnce executes a single HTTP request without retry.
func (h *httpClient) doOnce(ctx context.Context, opts requestOptions, dest interface{}) error {
	var bodyReader io.Reader
	if opts.body != nil {
		b, err := json.Marshal(opts.body)
		if err != nil {
			return &decodeError{err: fmt.Errorf("xcrop: failed to marshal request body: %w", err)}
		}
		bodyReader = bytes.NewReader(b)
	}

	// Build URL with properly encoded query parameters (#1)
	reqURL := h.baseURL + opts.path
	if len(opts.query) > 0 {
		params := url.Values{}
		for k, v := range opts.query {
			params.Set(k, v)
		}
		reqURL += "?" + params.Encode()
	}

	// Log warning if baseURL is not HTTPS (#10)
	if !strings.HasPrefix(h.baseURL, "https://") && !strings.HasPrefix(h.baseURL, "http://localhost") && !strings.HasPrefix(h.baseURL, "http://127.0.0.1") {
		log.Printf("xcrop: WARNING: base URL %q is not using HTTPS. API keys may be transmitted in plaintext.", h.baseURL)
	}

	req, err := http.NewRequestWithContext(ctx, opts.method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("xcrop: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("User-Agent", "xcrop-go/"+Version) // #13: use Version constant
	if opts.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("xcrop: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("xcrop: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		var errResp struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			apiErr.Message = errResp.Error
			apiErr.Code = errResp.Code
		} else {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		// Store Retry-After header for rate limit responses (#4)
		if resp.StatusCode == 429 {
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if seconds, err := strconv.Atoi(ra); err == nil {
					apiErr.RetryAfter = time.Duration(seconds) * time.Second
				}
			}
		}
		return apiErr
	}

	if dest != nil {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return &decodeError{err: fmt.Errorf("xcrop: failed to decode response: %w", err)}
		}
	}

	return nil
}

// calcBackoff computes the retry delay using exponential backoff with jitter (#14).
// For 429 responses, it respects the Retry-After header if present (#4).
func (h *httpClient) calcBackoff(attempt int, lastErr error) time.Duration {
	// Check for Retry-After from rate limit errors (#4)
	if apiErr, ok := lastErr.(*APIError); ok && apiErr.StatusCode == 429 {
		if apiErr.RetryAfter > 0 {
			// Add small jitter to Retry-After to avoid thundering herd
			jitter := time.Duration(rand.Int63n(int64(500 * time.Millisecond)))
			return apiErr.RetryAfter + jitter
		}
		// Use a longer base for rate limits when no Retry-After header
		base := time.Duration(2000) * time.Millisecond
		backoff := base * time.Duration(math.Pow(2, float64(attempt-1)))
		if backoff > time.Duration(defaultRetryMaxMs)*time.Millisecond {
			backoff = time.Duration(defaultRetryMaxMs) * time.Millisecond
		}
		// Add jitter: ±25% of backoff
		jitter := time.Duration(rand.Int63n(int64(backoff) / 2))
		return backoff/2 + jitter
	}

	// Standard exponential backoff for 5xx with jitter
	base := time.Duration(defaultRetryBaseMs) * time.Millisecond
	backoff := base * time.Duration(math.Pow(2, float64(attempt-1)))
	if backoff > time.Duration(defaultRetryMaxMs)*time.Millisecond {
		backoff = time.Duration(defaultRetryMaxMs) * time.Millisecond
	}
	// Add jitter: ±25% of backoff
	jitter := time.Duration(rand.Int63n(int64(backoff) / 2))
	return backoff/2 + jitter
}

// buildPaginationQuery adds count and cursor to a query map.
func buildPaginationQuery(params *PaginationParams) map[string]string {
	q := make(map[string]string)
	if params != nil {
		if params.Count > 0 {
			q["count"] = strconv.Itoa(params.Count)
		}
		if params.Cursor != "" {
			q["cursor"] = params.Cursor
		}
	}
	return q
}
