package xcrop

import (
	"context"
	"fmt"
)

// ListsService handles list-related API endpoints.
type ListsService struct {
	http *httpClient
}

func newListsService(h *httpClient) *ListsService {
	return &ListsService{http: h}
}

// ListTweets retrieves tweets from a list. Returns an iterator for automatic pagination.
func (s *ListsService) ListTweets(ctx context.Context, listID string, params *PaginationParams) *Iterator[Tweet] {
	if listID == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: listId must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/lists/"+listID+"/tweets", query, nil))
}

// GetTweets retrieves a single page of tweets from a list.
func (s *ListsService) GetTweets(ctx context.Context, listID string, params *PaginationParams) ([]Tweet, Meta, error) {
	if listID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: listId must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/lists/" + listID + "/tweets",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListMembers retrieves members of a list. Returns an iterator for automatic pagination.
func (s *ListsService) ListMembers(ctx context.Context, listID string, params *PaginationParams) *Iterator[User] {
	if listID == "" {
		return errIterator[User](fmt.Errorf("xcrop: listId must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/lists/"+listID+"/members", query, nil))
}

// GetMembers retrieves a single page of list members.
func (s *ListsService) GetMembers(ctx context.Context, listID string, params *PaginationParams) ([]User, Meta, error) {
	if listID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: listId must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/lists/" + listID + "/members",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListSubscribers retrieves subscribers of a list. Returns an iterator for automatic pagination.
func (s *ListsService) ListSubscribers(ctx context.Context, listID string, params *PaginationParams) *Iterator[User] {
	if listID == "" {
		return errIterator[User](fmt.Errorf("xcrop: listId must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/lists/"+listID+"/subscribers", query, nil))
}

// GetSubscribers retrieves a single page of list subscribers.
func (s *ListsService) GetSubscribers(ctx context.Context, listID string, params *PaginationParams) ([]User, Meta, error) {
	if listID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: listId must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/lists/" + listID + "/subscribers",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}
