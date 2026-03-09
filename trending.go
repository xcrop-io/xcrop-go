package xcrop

import (
	"context"
)

// TrendingService handles the trending topics endpoint.
type TrendingService struct {
	http *httpClient
}

func newTrendingService(h *httpClient) *TrendingService {
	return &TrendingService{http: h}
}

// TrendingResponse wraps trending topics with metadata.
type TrendingResponse struct {
	Data []TrendingTopic `json:"data"`
	Meta Meta            `json:"meta"`
}

// Get retrieves the current trending topics.
func (s *TrendingService) Get(ctx context.Context) ([]TrendingTopic, error) {
	var resp TrendingResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/trending",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
