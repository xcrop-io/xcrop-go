package xcrop

import (
	"errors"
	"fmt"
	"time"
)

// APIError represents an error returned by the XCROP API.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Message is the human-readable error message.
	Message string
	// Code is the machine-readable error code (e.g., "RATE_LIMITED", "NOT_FOUND").
	Code string
	// RetryAfter is the duration from the Retry-After header (429 responses only).
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("xcrop: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("xcrop: %s (status=%d)", e.Message, e.StatusCode)
}

// IsNotFound reports whether the error is a 404 Not Found.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsRateLimited reports whether the error is a 429 Too Many Requests.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429
	}
	return false
}

// IsUnauthorized reports whether the error is a 401 Unauthorized.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401
	}
	return false
}

// IsForbidden reports whether the error is a 403 Forbidden.
func IsForbidden(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 403
	}
	return false
}

// IsServerError reports whether the error is a 5xx server error.
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}
