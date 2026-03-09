package xcrop

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// KOLService handles the KOL (Key Opinion Leader) timeline endpoint.
type KOLService struct {
	http *httpClient
}

func newKOLService(h *httpClient) *KOLService {
	return &KOLService{http: h}
}

// KOLTimelineResponse wraps KOL timeline results with metadata.
type KOLTimelineResponse struct {
	Data []Tweet `json:"data"`
	Meta Meta    `json:"meta"`
}

// Timeline retrieves a merged timeline of tweets from multiple KOL users.
func (s *KOLService) Timeline(ctx context.Context, params *KOLTimelineParams) ([]Tweet, Meta, error) {
	if params == nil || len(params.Usernames) == 0 {
		return nil, Meta{}, fmt.Errorf("xcrop: KOLTimelineParams.Usernames must not be empty")
	}

	query := make(map[string]string)
	query["usernames"] = strings.Join(params.Usernames, ",")
	if params.Count > 0 {
		query["count"] = strconv.Itoa(params.Count)
	}

	var resp KOLTimelineResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/kol/timeline",
		query:  query,
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}
