package xcrop

import (
	"context"
	"fmt"
)

// UsersService handles user-related API endpoints.
type UsersService struct {
	http *httpClient
}

func newUsersService(h *httpClient) *UsersService {
	return &UsersService{http: h}
}

// GetResponse wraps a single user with metadata.
type GetUserResponse struct {
	Data User `json:"data"`
	Meta Meta `json:"meta"`
}

// Get retrieves a user profile by username.
func (s *UsersService) Get(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp GetUserResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// TweetsResponse wraps a list of tweets with metadata.
type TweetsResponse struct {
	Data []Tweet `json:"data"`
	Meta Meta    `json:"meta"`
}

// ListTweets retrieves tweets by a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListTweets(ctx context.Context, username string, params *PaginationParams) *Iterator[Tweet] {
	if username == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/users/"+username+"/tweets", query, nil))
}

// GetTweets retrieves a single page of tweets by a user.
func (s *UsersService) GetTweets(ctx context.Context, username string, params *PaginationParams) ([]Tweet, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/tweets",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListMentions retrieves mentions of a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListMentions(ctx context.Context, username string, params *PaginationParams) *Iterator[Tweet] {
	if username == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/users/"+username+"/mentions", query, nil))
}

// GetMentions retrieves a single page of mentions of a user.
func (s *UsersService) GetMentions(ctx context.Context, username string, params *PaginationParams) ([]Tweet, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/mentions",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// UsersResponse wraps a list of users with metadata.
type UsersResponse struct {
	Data []User `json:"data"`
	Meta Meta   `json:"meta"`
}

// ListFollowers retrieves followers of a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListFollowers(ctx context.Context, username string, params *PaginationParams) *Iterator[User] {
	if username == "" {
		return errIterator[User](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/users/"+username+"/followers", query, nil))
}

// GetFollowers retrieves a single page of followers.
func (s *UsersService) GetFollowers(ctx context.Context, username string, params *PaginationParams) ([]User, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/followers",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListFollowing retrieves users that a user follows. Returns an iterator for automatic pagination.
func (s *UsersService) ListFollowing(ctx context.Context, username string, params *PaginationParams) *Iterator[User] {
	if username == "" {
		return errIterator[User](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/users/"+username+"/following", query, nil))
}

// GetFollowing retrieves a single page of following.
func (s *UsersService) GetFollowing(ctx context.Context, username string, params *PaginationParams) ([]User, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/following",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListReplies retrieves replies by a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListReplies(ctx context.Context, username string, params *PaginationParams) *Iterator[Tweet] {
	if username == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/users/"+username+"/replies", query, nil))
}

// GetReplies retrieves a single page of replies by a user.
func (s *UsersService) GetReplies(ctx context.Context, username string, params *PaginationParams) ([]Tweet, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/replies",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListMedia retrieves media tweets by a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListMedia(ctx context.Context, username string, params *PaginationParams) *Iterator[Tweet] {
	if username == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/users/"+username+"/media", query, nil))
}

// GetMedia retrieves a single page of media tweets by a user.
func (s *UsersService) GetMedia(ctx context.Context, username string, params *PaginationParams) ([]Tweet, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/media",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListVerifiedFollowers retrieves blue-verified followers of a user. Returns an iterator.
func (s *UsersService) ListVerifiedFollowers(ctx context.Context, username string, params *PaginationParams) *Iterator[User] {
	if username == "" {
		return errIterator[User](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[User](s.http, "GET", "/users/"+username+"/verified-followers", query, nil))
}

// GetVerifiedFollowers retrieves a single page of verified followers.
func (s *UsersService) GetVerifiedFollowers(ctx context.Context, username string, params *PaginationParams) ([]User, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/verified-followers",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListLikes retrieves tweets liked by a user. Returns an iterator for automatic pagination.
func (s *UsersService) ListLikes(ctx context.Context, username string, params *PaginationParams) *Iterator[Tweet] {
	if username == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: username must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/users/"+username+"/likes", query, nil))
}

// GetLikes retrieves a single page of liked tweets.
func (s *UsersService) GetLikes(ctx context.Context, username string, params *PaginationParams) ([]Tweet, Meta, error) {
	if username == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/" + username + "/likes",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// BatchGet retrieves multiple users by usernames in a single request. Max 100.
func (s *UsersService) BatchGet(ctx context.Context, usernames []string) ([]User, error) {
	if len(usernames) == 0 {
		return nil, fmt.Errorf("xcrop: usernames must not be empty")
	}
	if len(usernames) > 100 {
		return nil, fmt.Errorf("xcrop: BatchGet supports max 100 usernames, got %d", len(usernames))
	}

	body := struct {
		Usernames []string `json:"usernames"`
	}{Usernames: usernames}

	var resp UsersResponse
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/users/batch",
		body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// RelationshipResponse wraps a relationship check result.
type RelationshipResponse struct {
	Data Relationship `json:"data"`
	Meta Meta         `json:"meta"`
}

// GetRelationship checks the follow relationship between two users.
func (s *UsersService) GetRelationship(ctx context.Context, source, target string) (*Relationship, error) {
	if source == "" {
		return nil, fmt.Errorf("xcrop: source username must not be empty")
	}
	if target == "" {
		return nil, fmt.Errorf("xcrop: target username must not be empty")
	}
	var resp RelationshipResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/users/relationship",
		query: map[string]string{
			"source": source,
			"target": target,
		},
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Follow follows a user. Requires a connected X account.
func (s *UsersService) Follow(ctx context.Context, username string) (*WriteResult, error) {
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/users/" + username + "/follow",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Unfollow unfollows a user. Requires a connected X account.
func (s *UsersService) Unfollow(ctx context.Context, username string) (*WriteResult, error) {
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "DELETE",
		path:   "/users/" + username + "/follow",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// errIterator returns an iterator that immediately yields an error.
func errIterator[T any](err error) *Iterator[T] {
	return &Iterator[T]{
		err:     err,
		started: true,
	}
}
