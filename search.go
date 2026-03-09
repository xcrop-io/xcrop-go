package xcrop

import (
	"context"
	"fmt"
)

// SearchService handles search API endpoints.
type SearchService struct {
	http *httpClient
}

func newSearchService(h *httpClient) *SearchService {
	return &SearchService{http: h}
}

// SearchTweetsResponse wraps search results with metadata.
type SearchTweetsResponse struct {
	Data []Tweet `json:"data"`
	Meta Meta    `json:"meta"`
}

// Tweets searches for tweets matching the given parameters.
func (s *SearchService) Tweets(ctx context.Context, params *SearchParams) ([]Tweet, Meta, error) {
	// Validate params (#7)
	if params == nil {
		return nil, Meta{}, fmt.Errorf("xcrop: SearchParams must not be nil")
	}
	if params.Query == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: SearchParams.Query must not be empty")
	}

	var resp SearchTweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/search",
		body:   params,
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// TweetsPaginate returns an iterator that automatically paginates through search results.
func (s *SearchService) TweetsPaginate(ctx context.Context, params *SearchParams) *Iterator[Tweet] {
	return newIterator(ctx, func(ctx context.Context, cursor string) ([]Tweet, Meta, error) {
		p := *params // copy to avoid mutating original
		if cursor != "" {
			p.Cursor = cursor
		}
		return s.Tweets(ctx, &p)
	})
}

// SearchUsersResponse wraps user search results with metadata.
type SearchUsersResponse struct {
	Data []User `json:"data"`
	Meta Meta   `json:"meta"`
}

// Users searches for users matching the given parameters.
func (s *SearchService) Users(ctx context.Context, params *UserSearchParams) ([]User, Meta, error) {
	// Validate params (#7)
	if params == nil {
		return nil, Meta{}, fmt.Errorf("xcrop: UserSearchParams must not be nil")
	}
	if params.Query == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: UserSearchParams.Query must not be empty")
	}

	var resp SearchUsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/search/users",
		body:   params,
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}
