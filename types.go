package xcrop

import "time"

// Meta contains response metadata returned by the API.
type Meta struct {
	LatencyMs int    `json:"latency_ms"`
	Cached    bool   `json:"cached"`
	Total     int    `json:"total,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	HasNext   bool   `json:"has_next,omitempty"`
}

// User represents an X/Twitter user profile.
type User struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	URL             string    `json:"url"`
	ProfileImageURL string    `json:"profile_image_url"`
	BannerURL       string    `json:"banner_url"`
	Verified        bool      `json:"verified"`
	IsBlueVerified  bool      `json:"is_blue_verified"`
	Protected       bool      `json:"protected"`
	FollowersCount  int       `json:"followers_count"`
	FollowingCount  int       `json:"following_count"`
	TweetCount      int       `json:"tweet_count"`
	LikesCount      int       `json:"likes_count"`
	ListedCount     int       `json:"listed_count"`
	CreatedAt       time.Time `json:"created_at"`
	PinnedTweetID   string    `json:"pinned_tweet_id,omitempty"`
}

// Tweet represents an X/Twitter tweet.
type Tweet struct {
	ID               string     `json:"id"`
	Text             string     `json:"text"`
	AuthorID         string     `json:"author_id"`
	Author           *User      `json:"author,omitempty"`
	ConversationID   string     `json:"conversation_id"`
	InReplyToUserID  string     `json:"in_reply_to_user_id,omitempty"`
	InReplyToTweetID string     `json:"in_reply_to_tweet_id,omitempty"`
	ReferencedTweets []RefTweet `json:"referenced_tweets,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	Language         string     `json:"language"`
	Source           string     `json:"source"`
	LikeCount        int        `json:"like_count"`
	RetweetCount     int        `json:"retweet_count"`
	ReplyCount       int        `json:"reply_count"`
	QuoteCount       int        `json:"quote_count"`
	ViewCount        int        `json:"view_count"`
	BookmarkCount    int        `json:"bookmark_count"`
	IsRetweet        bool       `json:"is_retweet"`
	IsReply          bool       `json:"is_reply"`
	IsQuote          bool       `json:"is_quote"`
	Entities         *Entities  `json:"entities,omitempty"`
	Media            []Media    `json:"media,omitempty"`
	QuotedTweet      *Tweet     `json:"quoted_tweet,omitempty"`
}

// RefTweet represents a referenced tweet (replied_to, quoted, retweeted).
type RefTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Entities contains tweet entity annotations.
type Entities struct {
	URLs     []EntityURL     `json:"urls,omitempty"`
	Mentions []EntityMention `json:"mentions,omitempty"`
	Hashtags []EntityHashtag `json:"hashtags,omitempty"`
}

// EntityURL represents a URL entity in a tweet.
type EntityURL struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

// EntityMention represents a user mention entity.
type EntityMention struct {
	Username string `json:"username"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

// EntityHashtag represents a hashtag entity.
type EntityHashtag struct {
	Tag   string `json:"tag"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Media represents a media attachment on a tweet.
type Media struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	PreviewURL string `json:"preview_url"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	DurationMs int    `json:"duration_ms,omitempty"`
	AltText    string `json:"alt_text,omitempty"`
}

// TrendingTopic represents a trending topic.
type TrendingTopic struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	TweetCount    int    `json:"tweet_count"`
	DomainContext string `json:"domain_context,omitempty"`
}

// Relationship represents the follow relationship between two users.
type Relationship struct {
	SourceFollowsTarget bool `json:"source_follows_target"`
	TargetFollowsSource bool `json:"target_follows_source"`
	Blocking            bool `json:"blocking"`
	Muting              bool `json:"muting"`
}

// AccountStatus represents the connected X account status.
type AccountStatus struct {
	Connected bool   `json:"connected"`
	Username  string `json:"username,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

// InteractionCheck represents the result of an interaction check.
type InteractionCheck struct {
	Found    bool   `json:"found"`
	TweetID  string `json:"tweet_id,omitempty"`
	Username string `json:"username"`
}

// List represents an X/Twitter list.
type List struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	MemberCount   int    `json:"member_count"`
	FollowerCount int    `json:"follower_count"`
	Private       bool   `json:"private"`
	Owner         *User  `json:"owner,omitempty"`
}

// WriteResult represents the result of a write operation.
type WriteResult struct {
	Success bool   `json:"success"`
	TweetID string `json:"tweet_id,omitempty"`
	Message string `json:"message,omitempty"`
}

// PaginationParams specifies pagination options for list endpoints.
type PaginationParams struct {
	Count  int
	Cursor string
}

// SearchParams specifies parameters for tweet search.
type SearchParams struct {
	Query           string `json:"query"`
	Count           int    `json:"count,omitempty"`
	Sort            string `json:"sort,omitempty"`
	Cursor          string `json:"cursor,omitempty"`
	Lang            string `json:"lang,omitempty"`
	MinLikes        int    `json:"min_likes,omitempty"`
	MinRetweets     int    `json:"min_retweets,omitempty"`
	ExcludeReplies  bool   `json:"exclude_replies,omitempty"`
	ExcludeRetweets bool   `json:"exclude_retweets,omitempty"`
	Since           string `json:"since,omitempty"`
	Until           string `json:"until,omitempty"`
}

// UserSearchParams specifies parameters for user search.
type UserSearchParams struct {
	Query string `json:"query"`
	Count int    `json:"count,omitempty"`
}

// ConnectParams specifies parameters for connecting an X account.
type ConnectParams struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	TOTPSecret string `json:"totp_secret,omitempty"`
}

// String returns a string representation of ConnectParams with sensitive fields redacted. (#9)
func (p ConnectParams) String() string {
	totp := ""
	if p.TOTPSecret != "" {
		totp = ", TOTPSecret: [REDACTED]"
	}
	return "ConnectParams{Username: " + p.Username + ", Password: [REDACTED]" + totp + "}"
}

// GoString returns a Go-syntax representation with sensitive fields redacted. (#9)
func (p ConnectParams) GoString() string {
	totp := ""
	if p.TOTPSecret != "" {
		totp = ", TOTPSecret: \"[REDACTED]\""
	}
	return "xcrop.ConnectParams{Username: \"" + p.Username + "\", Password: \"[REDACTED]\"" + totp + "}"
}

// KOLTimelineParams specifies parameters for the KOL timeline endpoint.
type KOLTimelineParams struct {
	Usernames []string
	Count     int
}
