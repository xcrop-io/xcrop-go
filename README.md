# xcrop-go

Go client library for the [XCROP API](https://xcrop.io) — X/Twitter data intelligence platform.

[![Go Reference](https://pkg.go.dev/badge/github.com/xcrop-io/xcrop-go.svg)](https://pkg.go.dev/github.com/xcrop-io/xcrop-go)

## Installation

```bash
go get github.com/xcrop-io/xcrop-go
```

Requires Go 1.21+. Zero external dependencies — uses only the standard library.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/xcrop-io/xcrop-go"
)

func main() {
    client := xcrop.NewClient("xc_live_your_api_key")
    ctx := context.Background()

    // Get a user profile
    user, err := client.Users.Get(ctx, "elonmusk")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s has %d followers\n", user.Name, user.FollowersCount)
}
```

## Features

- **41 API endpoints** — users, tweets, search, lists, trending, KOL timeline, write operations, interaction checks
- **Auto-retry** — exponential backoff on 429 (rate limit) and 5xx (server error), max 3 retries
- **Pagination iterator** — iterate through pages automatically with `Next()`/`Item()` pattern
- **Context support** — all methods accept `context.Context` for cancellation and timeouts
- **Configurable** — custom HTTP client, base URL, timeout, max retries
- **Typed errors** — `*xcrop.APIError` with status code, message, and error code
- **Zero dependencies** — only Go standard library

## Usage

### Client Options

```go
// Default client
client := xcrop.NewClient("xc_live_...")

// With options
client := xcrop.NewClient("xc_live_...",
    xcrop.WithTimeout(10 * time.Second),
    xcrop.WithMaxRetries(5),
    xcrop.WithBaseURL("https://custom.api.com/v2"),
    xcrop.WithHTTPClient(&http.Client{Transport: customTransport}),
)
```

### Users

```go
// Get user profile
user, err := client.Users.Get(ctx, "elonmusk")

// Get user tweets (single page)
tweets, meta, err := client.Users.GetTweets(ctx, "elonmusk", &xcrop.PaginationParams{
    Count: 50,
})

// Get followers
followers, meta, err := client.Users.GetFollowers(ctx, "elonmusk", &xcrop.PaginationParams{
    Count: 100,
})

// Batch lookup (max 100 usernames)
users, err := client.Users.BatchGet(ctx, []string{"elonmusk", "jack", "vaborCFA"})

// Check relationship
rel, err := client.Users.GetRelationship(ctx, "elonmusk", "jack")
fmt.Println(rel.SourceFollowsTarget) // true/false

// Other user endpoints
mentions, meta, err := client.Users.GetMentions(ctx, "elonmusk", params)
replies, meta, err := client.Users.GetReplies(ctx, "elonmusk", params)
media, meta, err := client.Users.GetMedia(ctx, "elonmusk", params)
verified, meta, err := client.Users.GetVerifiedFollowers(ctx, "elonmusk", params)
likes, meta, err := client.Users.GetLikes(ctx, "elonmusk", params)
following, meta, err := client.Users.GetFollowing(ctx, "elonmusk", params)
```

### Tweets

```go
// Get a single tweet
tweet, err := client.Tweets.Get(ctx, "1234567890")

// Get conversation thread
replies, meta, err := client.Tweets.GetConversation(ctx, "1234567890", &xcrop.PaginationParams{
    Count: 50,
})

// Get quote tweets, likers, retweeters
quotes, meta, err := client.Tweets.GetQuotes(ctx, "1234567890", params)
likers, meta, err := client.Tweets.GetLikers(ctx, "1234567890", params)
retweeters, meta, err := client.Tweets.GetRetweeters(ctx, "1234567890", params)

// Batch lookup (max 100 IDs)
tweets, err := client.Tweets.BatchGet(ctx, []string{"123", "456", "789"})
```

### Search

```go
// Search tweets
tweets, meta, err := client.Search.Tweets(ctx, &xcrop.SearchParams{
    Query:           "bitcoin",
    Count:           50,
    Sort:            "popular",  // "latest" or "popular"
    Lang:            "en",
    MinLikes:        100,
    MinRetweets:     10,
    ExcludeReplies:  true,
    ExcludeRetweets: true,
    Since:           "2024-01-01",
    Until:           "2024-12-31",
})

// Search users
users, meta, err := client.Search.Users(ctx, &xcrop.UserSearchParams{
    Query: "crypto",
    Count: 20,
})

// Paginate through search results
iter := client.Search.TweetsPaginate(ctx, &xcrop.SearchParams{
    Query: "ethereum",
    Count: 100,
})
for iter.Next() {
    fmt.Println(iter.Item().Text)
}
```

### Pagination

All list endpoints have both single-page (`Get*`) and iterator (`List*`) variants.

```go
// Iterator pattern — automatically fetches next pages
iter := client.Users.ListTweets(ctx, "elonmusk", &xcrop.PaginationParams{Count: 100})
for iter.Next() {
    tweet := iter.Item()
    fmt.Println(tweet.Text)
}
if err := iter.Err(); err != nil {
    log.Fatal(err)
}

// Collect all items at once (caution with large datasets)
allTweets, err := client.Users.ListTweets(ctx, "elonmusk", &xcrop.PaginationParams{
    Count: 100,
}).Collect()

// Manual pagination with cursors
tweets, meta, err := client.Users.GetTweets(ctx, "elonmusk", &xcrop.PaginationParams{
    Count: 50,
})
// Next page:
if meta.HasNext {
    nextPage, meta, err := client.Users.GetTweets(ctx, "elonmusk", &xcrop.PaginationParams{
        Count:  50,
        Cursor: meta.Cursor,
    })
}
```

### Lists

```go
tweets, meta, err := client.Lists.GetTweets(ctx, "1234567890", params)
members, meta, err := client.Lists.GetMembers(ctx, "1234567890", params)
subscribers, meta, err := client.Lists.GetSubscribers(ctx, "1234567890", params)
```

### Trending & KOL

```go
// Trending topics
topics, err := client.Trending.Get(ctx)
for _, t := range topics {
    fmt.Printf("%s (%d tweets)\n", t.Name, t.TweetCount)
}

// KOL timeline (merged feed from multiple users)
tweets, meta, err := client.KOL.Timeline(ctx, &xcrop.KOLTimelineParams{
    Usernames: []string{"elonmusk", "vaborCFA", "CryptoCapo_"},
    Count:     50,
})
```

### Write Operations

Requires a connected X account via `Account.Connect()`.

```go
// Connect your X account
_, err := client.Account.Connect(ctx, &xcrop.ConnectParams{
    Username:   "your_username",
    Password:   "your_password",
    TOTPSecret: "OPTIONAL_2FA_SECRET",
})

// Check connection status
status, err := client.Account.Status(ctx)
fmt.Println(status.Connected, status.Username)

// Create a tweet
result, err := client.Tweets.Create(ctx, "Hello from XCROP Go SDK!")
fmt.Println(result.TweetID)

// Reply, quote, delete
client.Tweets.Reply(ctx, "1234567890", "Great point!")
client.Tweets.Quote(ctx, "1234567890", "This is interesting")
client.Tweets.Delete(ctx, result.TweetID)

// Like, retweet
client.Tweets.Like(ctx, "1234567890")
client.Tweets.Unlike(ctx, "1234567890")
client.Tweets.Retweet(ctx, "1234567890")
client.Tweets.Unretweet(ctx, "1234567890")

// Follow, unfollow
client.Users.Follow(ctx, "elonmusk")
client.Users.Unfollow(ctx, "elonmusk")

// Disconnect when done
client.Account.Disconnect(ctx)
```

### Interaction Checks

Check if a user performed specific interactions on a tweet. Does not require a connected account.

```go
// Check if user retweeted
check, err := client.Tweets.CheckRetweet(ctx, "1234567890", "elonmusk")
fmt.Println(check.Found) // true or false

// Check reply, quote
client.Tweets.CheckReply(ctx, "1234567890", "elonmusk")
client.Tweets.CheckQuote(ctx, "1234567890", "elonmusk")

// Check like (note: X has hidden likes, may be unavailable)
client.Tweets.CheckLike(ctx, "1234567890", "elonmusk")
```

### Error Handling

```go
user, err := client.Users.Get(ctx, "nonexistent_user")
if err != nil {
    if xcrop.IsNotFound(err) {
        fmt.Println("User not found")
    } else if xcrop.IsRateLimited(err) {
        fmt.Println("Rate limited (retries exhausted)")
    } else if xcrop.IsUnauthorized(err) {
        fmt.Println("Invalid API key")
    } else if apiErr, ok := err.(*xcrop.APIError); ok {
        fmt.Printf("API error: %s (code=%s, status=%d)\n",
            apiErr.Message, apiErr.Code, apiErr.StatusCode)
    } else {
        fmt.Printf("Network/other error: %v\n", err)
    }
}
```

### Context & Cancellation

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := client.Users.Get(ctx, "elonmusk")

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel all in-flight requests
}()

iter := client.Users.ListFollowers(ctx, "elonmusk", &xcrop.PaginationParams{Count: 200})
for iter.Next() {
    // Will stop when context is cancelled
}
```

## API Reference

Full API documentation: [https://xcrop.io/docs](https://xcrop.io/docs)

### Services

| Service | Description |
|---------|-------------|
| `client.Users` | User profiles, tweets, followers, following, mentions, replies, media, likes, batch, relationship, follow/unfollow |
| `client.Tweets` | Single tweet, conversation, quotes, likers, retweeters, batch, create/reply/quote/delete, like/unlike, retweet/unretweet, interaction checks |
| `client.Search` | Tweet search, user search |
| `client.Lists` | List tweets, members, subscribers |
| `client.Trending` | Trending topics |
| `client.KOL` | KOL merged timeline |
| `client.Account` | Connect/disconnect X account, check status |
| `client.Stream` | SSE real-time stream |

### Real-time Stream (SSE)

```go
// Connect to real-time stream
resp, err := client.Stream.Connect(ctx, &xcrop.StreamParams{
    Keywords:  []string{"bitcoin", "ethereum"},
    Usernames: []string{"elonmusk"},
})
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

scanner := bufio.NewScanner(resp.Body)
for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "data: ") {
        fmt.Println(line[6:])
    }
}
```

## License

MIT
