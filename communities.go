package xcrop

import (
	"context"
	"fmt"
	"strconv"
)

// CommunitiesService handles X Community endpoints.
//
// Beta: these endpoints are under active development on the API and currently
// return a 503 Service Unavailable while the backend integration is finished.
// Use IsServerError(err) or inspect (*APIError).StatusCode to detect this.
type CommunitiesService struct {
	http *httpClient
}

func newCommunitiesService(h *httpClient) *CommunitiesService {
	return &CommunitiesService{http: h}
}

// CommunityResponse wraps a single community with metadata.
type CommunityResponse struct {
	Data Community `json:"data"`
	Meta Meta      `json:"meta"`
}

// Get retrieves community details by ID.
//
// Beta: may currently return a 503 while this endpoint is under development.
func (s *CommunitiesService) Get(ctx context.Context, communityID string) (*Community, error) {
	if communityID == "" {
		return nil, fmt.Errorf("xcrop: communityId must not be empty")
	}
	var resp CommunityResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/communities/" + communityID,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CommunityListParams specifies pagination and sort options for community sub-resources.
type CommunityListParams struct {
	Count  int
	Cursor string
	// Sort accepts "latest", "popular", or "engagement" for ListTweets/GetTweets,
	// and "default", "followers", or "name" for ListMembers/GetMembers.
	Sort string
}

func buildCommunityQuery(params *CommunityListParams) map[string]string {
	q := make(map[string]string)
	if params != nil {
		if params.Count > 0 {
			q["count"] = strconv.Itoa(params.Count)
		}
		if params.Cursor != "" {
			q["cursor"] = params.Cursor
		}
		if params.Sort != "" {
			q["sort"] = params.Sort
		}
	}
	return q
}

// ListTweets retrieves tweets from a community timeline. Returns an iterator for
// automatic pagination.
//
// Beta: may currently return a 503 while this endpoint is under development.
func (s *CommunitiesService) ListTweets(ctx context.Context, communityID string, params *CommunityListParams) *Iterator[Tweet] {
	if communityID == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: communityId must not be empty"))
	}
	query := buildCommunityQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/communities/"+communityID+"/tweets", query, nil))
}

// GetTweets retrieves a single page of tweets from a community timeline.
//
// Beta: may currently return a 503 while this endpoint is under development.
func (s *CommunitiesService) GetTweets(ctx context.Context, communityID string, params *CommunityListParams) ([]Tweet, Meta, error) {
	if communityID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: communityId must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/communities/" + communityID + "/tweets",
		query:  buildCommunityQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListMembers retrieves members of a community. Returns an iterator for automatic
// pagination.
//
// Beta: may currently return a 503 while this endpoint is under development.
func (s *CommunitiesService) ListMembers(ctx context.Context, communityID string, params *CommunityListParams) *Iterator[User] {
	if communityID == "" {
		return errIterator[User](fmt.Errorf("xcrop: communityId must not be empty"))
	}
	query := buildCommunityQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/communities/"+communityID+"/members", query, nil))
}

// GetMembers retrieves a single page of community members.
//
// Beta: may currently return a 503 while this endpoint is under development.
func (s *CommunitiesService) GetMembers(ctx context.Context, communityID string, params *CommunityListParams) ([]User, Meta, error) {
	if communityID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: communityId must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/communities/" + communityID + "/members",
		query:  buildCommunityQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}
