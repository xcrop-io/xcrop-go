package xcrop

import (
	"context"
	"fmt"
)

// TweetsService handles tweet-related API endpoints.
type TweetsService struct {
	http *httpClient
}

func newTweetsService(h *httpClient) *TweetsService {
	return &TweetsService{http: h}
}

// GetTweetResponse wraps a single tweet with metadata.
type GetTweetResponse struct {
	Data Tweet `json:"data"`
	Meta Meta  `json:"meta"`
}

// Get retrieves a single tweet by ID.
func (s *TweetsService) Get(ctx context.Context, tweetID string) (*Tweet, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp GetTweetResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListConversation retrieves the conversation thread for a tweet. Returns an iterator.
func (s *TweetsService) ListConversation(ctx context.Context, tweetID string, params *PaginationParams) *Iterator[Tweet] {
	if tweetID == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: tweetId must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/tweets/"+tweetID+"/conversation", query, nil))
}

// GetConversation retrieves a single page of a tweet's conversation.
func (s *TweetsService) GetConversation(ctx context.Context, tweetID string, params *PaginationParams) ([]Tweet, Meta, error) {
	if tweetID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID + "/conversation",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// ListQuotes retrieves quote tweets of a tweet. Returns an iterator.
func (s *TweetsService) ListQuotes(ctx context.Context, tweetID string, params *PaginationParams) *Iterator[Tweet] {
	if tweetID == "" {
		return errIterator[Tweet](fmt.Errorf("xcrop: tweetId must not be empty"))
	}
	query := buildPaginationQuery(params)
	return newIterator(ctx, makePaginatedFetcher[Tweet](s.http, "GET", "/tweets/"+tweetID+"/quotes", query, nil))
}

// GetQuotes retrieves a single page of quote tweets.
func (s *TweetsService) GetQuotes(ctx context.Context, tweetID string, params *PaginationParams) ([]Tweet, Meta, error) {
	if tweetID == "" {
		return nil, Meta{}, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID + "/quotes",
		query:  buildPaginationQuery(params),
	}, &resp)
	if err != nil {
		return nil, Meta{}, err
	}
	return resp.Data, resp.Meta, nil
}

// BatchGet retrieves multiple tweets by IDs in a single request. Max 100.
func (s *TweetsService) BatchGet(ctx context.Context, tweetIDs []string) ([]Tweet, error) {
	if len(tweetIDs) == 0 {
		return nil, fmt.Errorf("xcrop: tweetIDs must not be empty")
	}
	if len(tweetIDs) > 100 {
		return nil, fmt.Errorf("xcrop: BatchGet supports max 100 tweet IDs, got %d", len(tweetIDs))
	}

	body := struct {
		TweetIDs []string `json:"tweet_ids"`
	}{TweetIDs: tweetIDs}

	var resp TweetsResponse
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/batch",
		body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Create creates a new tweet. Requires a connected X account.
func (s *TweetsService) Create(ctx context.Context, text string) (*WriteResult, error) {
	if text == "" {
		return nil, fmt.Errorf("xcrop: tweet text must not be empty")
	}
	body := struct {
		Text string `json:"text"`
	}{Text: text}

	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/create",
		body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Reply replies to a tweet. Requires a connected X account.
func (s *TweetsService) Reply(ctx context.Context, tweetID, text string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	if text == "" {
		return nil, fmt.Errorf("xcrop: reply text must not be empty")
	}
	body := struct {
		TweetID string `json:"tweet_id"`
		Text    string `json:"text"`
	}{TweetID: tweetID, Text: text}

	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/reply",
		body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Quote creates a quote tweet. Requires a connected X account.
func (s *TweetsService) Quote(ctx context.Context, tweetID, text string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	if text == "" {
		return nil, fmt.Errorf("xcrop: quote text must not be empty")
	}
	body := struct {
		TweetID string `json:"tweet_id"`
		Text    string `json:"text"`
	}{TweetID: tweetID, Text: text}

	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/quote",
		body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a tweet. Requires a connected X account.
func (s *TweetsService) Delete(ctx context.Context, tweetID string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "DELETE",
		path:   "/tweets/" + tweetID,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Like likes a tweet. Requires a connected X account.
func (s *TweetsService) Like(ctx context.Context, tweetID string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/" + tweetID + "/like",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Unlike unlikes a tweet. Requires a connected X account.
func (s *TweetsService) Unlike(ctx context.Context, tweetID string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "DELETE",
		path:   "/tweets/" + tweetID + "/like",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Retweet retweets a tweet. Requires a connected X account.
func (s *TweetsService) Retweet(ctx context.Context, tweetID string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "POST",
		path:   "/tweets/" + tweetID + "/retweet",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Unretweet removes a retweet. Requires a connected X account.
func (s *TweetsService) Unretweet(ctx context.Context, tweetID string) (*WriteResult, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	var resp struct {
		Data WriteResult `json:"data"`
		Meta Meta        `json:"meta"`
	}
	err := s.http.do(ctx, requestOptions{
		method: "DELETE",
		path:   "/tweets/" + tweetID + "/retweet",
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// InteractionCheckResponse wraps an interaction check result.
type InteractionCheckResponse struct {
	Data InteractionCheck `json:"data"`
	Meta Meta             `json:"meta"`
}

// CheckRetweet checks if a user retweeted a tweet.
func (s *TweetsService) CheckRetweet(ctx context.Context, tweetID, username string) (*InteractionCheck, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp InteractionCheckResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID + "/check-retweet",
		query:  map[string]string{"username": username},
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CheckReply checks if a user replied to a tweet.
func (s *TweetsService) CheckReply(ctx context.Context, tweetID, username string) (*InteractionCheck, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp InteractionCheckResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID + "/check-reply",
		query:  map[string]string{"username": username},
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CheckQuote checks if a user quote-tweeted a tweet.
func (s *TweetsService) CheckQuote(ctx context.Context, tweetID, username string) (*InteractionCheck, error) {
	if tweetID == "" {
		return nil, fmt.Errorf("xcrop: tweetId must not be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("xcrop: username must not be empty")
	}
	var resp InteractionCheckResponse
	err := s.http.do(ctx, requestOptions{
		method: "GET",
		path:   "/tweets/" + tweetID + "/check-quote",
		query:  map[string]string{"username": username},
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
