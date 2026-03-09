package xcrop

import (
	"context"
	"encoding/json"
)

// Iterator provides a convenient way to iterate over paginated results.
// Use Next() to advance and Item() to access the current item.
//
// Example:
//
//	iter := client.Users.ListTweets(ctx, "elonmusk", &PaginationParams{Count: 100})
//	for iter.Next() {
//	    tweet := iter.Item()
//	    fmt.Println(tweet.Text)
//	}
//	if err := iter.Err(); err != nil {
//	    log.Fatal(err)
//	}
type Iterator[T any] struct {
	ctx     context.Context
	fetch   func(ctx context.Context, cursor string) ([]T, Meta, error)
	items   []T
	index   int
	cursor  string
	hasNext bool
	started bool
	err     error
	meta    Meta
}

// newIterator creates a new paginated iterator.
// fetch is called to load each page; it receives a cursor (empty for the first page)
// and returns items, metadata, and any error.
func newIterator[T any](ctx context.Context, fetch func(ctx context.Context, cursor string) ([]T, Meta, error)) *Iterator[T] {
	return &Iterator[T]{
		ctx:     ctx,
		fetch:   fetch,
		hasNext: true,
	}
}

// Next advances the iterator to the next item. It returns false when iteration
// is complete or an error has occurred. Check Err() after the loop.
func (it *Iterator[T]) Next() bool {
	if it.err != nil {
		return false
	}

	// Try to advance within the current page
	if it.started && it.index < len(it.items)-1 {
		it.index++
		return true
	}

	// Need to fetch next page
	if it.started && !it.hasNext {
		return false
	}

	// Check context
	if it.ctx.Err() != nil {
		it.err = it.ctx.Err()
		return false
	}

	items, meta, err := it.fetch(it.ctx, it.cursor)
	if err != nil {
		it.err = err
		return false
	}

	it.meta = meta
	it.cursor = meta.Cursor
	it.hasNext = meta.HasNext
	it.items = items
	it.index = 0
	it.started = true

	return len(items) > 0
}

// Item returns the current item. Must be called after a successful Next().
// Returns the zero value of T if called without a preceding successful Next(). (#5)
func (it *Iterator[T]) Item() T {
	if len(it.items) == 0 || it.index < 0 || it.index >= len(it.items) {
		var zero T
		return zero
	}
	return it.items[it.index]
}

// Items returns all items from the current page.
func (it *Iterator[T]) Items() []T {
	return it.items
}

// Err returns the first error encountered during iteration, if any.
func (it *Iterator[T]) Err() error {
	return it.err
}

// Meta returns the metadata from the most recent API response.
func (it *Iterator[T]) Meta() Meta {
	return it.meta
}

// Collect fetches all remaining pages and returns all items.
// This may make multiple API calls and should be used carefully with large datasets.
func (it *Iterator[T]) Collect() ([]T, error) {
	var all []T
	for it.Next() {
		all = append(all, it.Item())
	}
	if it.err != nil {
		return nil, it.err
	}
	return all, nil
}

// paginatedResponse is used to decode paginated API responses generically.
type paginatedResponse struct {
	Data json.RawMessage `json:"data"`
	Meta Meta            `json:"meta"`
}

// makePaginatedFetcher creates a fetch function for use with newIterator.
// It handles building the request, decoding the response, and extracting items.
func makePaginatedFetcher[T any](h *httpClient, method, path string, baseQuery map[string]string, body interface{}) func(ctx context.Context, cursor string) ([]T, Meta, error) {
	return func(ctx context.Context, cursor string) ([]T, Meta, error) {
		query := make(map[string]string)
		for k, v := range baseQuery {
			query[k] = v
		}
		if cursor != "" {
			query["cursor"] = cursor
		}

		var resp paginatedResponse
		err := h.do(ctx, requestOptions{
			method: method,
			path:   path,
			query:  query,
			body:   body,
		}, &resp)
		if err != nil {
			return nil, Meta{}, err
		}

		// Handle null data field by returning empty slice instead of error (#12)
		if len(resp.Data) == 0 || string(resp.Data) == "null" {
			return []T{}, resp.Meta, nil
		}

		var items []T
		if err := json.Unmarshal(resp.Data, &items); err != nil {
			// Try single item (some endpoints return object instead of array)
			var single T
			if err2 := json.Unmarshal(resp.Data, &single); err2 == nil {
				items = []T{single}
			} else {
				return nil, Meta{}, err
			}
		}

		return items, resp.Meta, nil
	}
}
