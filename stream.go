package xcrop

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// StreamService handles the SSE real-time stream endpoint.
type StreamService struct {
	h *httpClient
}

func newStreamService(h *httpClient) *StreamService {
	return &StreamService{h: h}
}

// StreamParams specifies parameters for the real-time stream endpoint.
type StreamParams struct {
	// Keywords to filter the stream (optional).
	Keywords []string
	// Usernames to filter the stream (optional).
	Usernames []string
}

// Connect opens a raw SSE connection to the /stream endpoint.
// The caller is responsible for reading from the response body and closing it
// when done. The response will have Content-Type "text/event-stream".
//
// Example:
//
//	resp, err := client.Stream.Connect(ctx, &xcrop.StreamParams{
//	    Keywords: []string{"bitcoin", "ethereum"},
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	scanner := bufio.NewScanner(resp.Body)
//	for scanner.Scan() {
//	    line := scanner.Text()
//	    // Parse SSE events: lines starting with "data: "
//	    if strings.HasPrefix(line, "data: ") {
//	        fmt.Println(line[6:])
//	    }
//	}
func (s *StreamService) Connect(ctx context.Context, params *StreamParams) (*http.Response, error) {
	reqURL := s.h.baseURL + "/stream"

	if params != nil {
		q := url.Values{}
		for _, kw := range params.Keywords {
			q.Add("keywords", kw)
		}
		for _, u := range params.Usernames {
			q.Add("usernames", u)
		}
		if encoded := q.Encode(); encoded != "" {
			reqURL += "?" + encoded
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("xcrop: failed to create stream request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.h.apiKey)
	req.Header.Set("User-Agent", "xcrop-go/"+Version)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := s.h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("xcrop: stream connection failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	return resp, nil
}
