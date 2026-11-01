package raims

import (
	"context"
	"encoding/json"
	"iter"
)

type pageInfo struct {
	HasNext        bool
	HasPrevious    bool
	NextAfter      *string
	PreviousBefore *string
}
type pageRequest struct {
	client    *Client
	key       string
	path      map[string]string
	params    map[string]json.RawMessage
	options   []RequestOption
	direction string
	seen      map[string]bool
}

func (r *Resource) setPageRequest(request *pageRequest) { r.pageRequest = request }

func clonePath(path map[string]string) map[string]string {
	result := map[string]string{}
	for key, value := range path {
		result[key] = value
	}
	return result
}
func cloneParams(params map[string]json.RawMessage) map[string]json.RawMessage {
	result := map[string]json.RawMessage{}
	for key, value := range params {
		result[key] = append(json.RawMessage(nil), value...)
	}
	return result
}

func fetchPage[T any](ctx context.Context, resource *Resource, page pageInfo, previous bool) (*T, error) {
	hasPage, cursor, field, opposite := page.HasNext, page.NextAfter, "after", "before"
	if previous {
		hasPage, cursor, field, opposite = page.HasPrevious, page.PreviousBefore, "before", "after"
	}
	if !hasPage {
		return nil, nil
	}
	if cursor == nil || *cursor == "" {
		return nil, &ProtocolError{Reason: "pagination continuation has no cursor", Response: resource.LastResponse}
	}
	if resource.pageRequest == nil {
		return nil, &ValidationError{Field: "page", Reason: "was not obtained from an SDK request"}
	}
	origin := resource.pageRequest
	params := cloneParams(origin.params)
	seen := map[string]bool{}
	if origin.direction == field {
		for value := range origin.seen {
			seen[value] = true
		}
	}
	var current string
	_ = json.Unmarshal(params[field], &current)
	if current != "" {
		seen[current] = true
	}
	if seen[*cursor] {
		return nil, &ProtocolError{Reason: "pagination cursor repeated", Response: resource.LastResponse}
	}
	seen[*cursor] = true
	params[field], _ = json.Marshal(*cursor)
	delete(params, opposite)
	result, err := request[T](ctx, origin.client, origin.key, origin.path, params, origin.options...)
	if err != nil {
		return nil, err
	}
	if collection, ok := any(result).(interface {
		pagination() pageInfo
		setPageRequest(*pageRequest)
	}); ok {
		collection.setPageRequest(&pageRequest{client: origin.client, key: origin.key, path: clonePath(origin.path), params: params, options: append([]RequestOption(nil), origin.options...), direction: field, seen: seen})
	}
	return result, nil
}

// Iterator loads pages only when Next or All advances. Do not share an iterator between goroutines.
type Iterator[T any] struct {
	ctx     context.Context
	client  *Client
	key     string
	path    map[string]string
	params  map[string]json.RawMessage
	options []RequestOption
	items   []*T
	current *T
	done    bool
	err     error
	seen    map[string]bool
}

func newIterator[T any](ctx context.Context, client *Client, key string, path map[string]string, params any, options ...RequestOption) *Iterator[T] {
	fields, err := marshalParams(params)
	seen := map[string]bool{}
	var initial string
	if json.Unmarshal(fields["after"], &initial) == nil && initial != "" {
		seen[initial] = true
	}
	return &Iterator[T]{ctx: ctx, client: client, key: key, path: clonePath(path), params: fields, options: append([]RequestOption(nil), options...), err: err, seen: seen}
}

// Next advances to the next item, fetching another page only when necessary.
// It returns false at the end or after an error; inspect Err to distinguish them.
func (i *Iterator[T]) Next() bool {
	i.current = nil
	for i.err == nil {
		if i.ctx == nil {
			i.err = &ValidationError{Field: "context", Reason: "must not be nil"}
			return false
		}
		if err := i.ctx.Err(); err != nil {
			i.err = err
			return false
		}
		if len(i.items) > 0 {
			i.current, i.items = i.items[0], i.items[1:]
			return true
		}
		if i.done {
			return false
		}
		var page struct {
			List []*T         `json:"list"`
			Page pageInfoJSON `json:"page"`
		}
		// Use the same schema validation and exact-name hydration as List.
		metadata, err := i.client.execute(i.ctx, i.key, i.path, i.params, &page, i.options...)
		if err != nil {
			i.err = err
			return false
		}
		i.items = page.List
		for _, item := range i.items {
			if resource, ok := any(item).(interface{ setResponse(*Response) }); ok {
				resource.setResponse(metadata)
			}
		}
		i.done = !page.Page.HasNext
		if page.Page.HasNext {
			cursor := page.Page.NextAfter
			if cursor == nil || *cursor == "" || i.seen[*cursor] {
				i.err = &ProtocolError{Reason: "pagination cursor is missing or repeated", Response: metadata}
				return false
			}
			i.seen[*cursor] = true
			i.params["after"], _ = json.Marshal(*cursor)
			delete(i.params, "before")
		}
	}
	return false
}

type pageInfoJSON struct {
	PageSize       int64   `json:"page_size"`
	HasNext        bool    `json:"has_next"`
	HasPrevious    bool    `json:"has_previous"`
	NextAfter      *string `json:"next_after"`
	PreviousBefore *string `json:"previous_before"`
}

// Current returns the item selected by the last successful Next call.
// It returns nil before iteration and after Next returns false.
func (i *Iterator[T]) Current() *T { return i.current }

// Err returns the error that stopped iteration, or nil when no error occurred.
func (i *Iterator[T]) Err() error { return i.err }

// All supports range over a lazy sequence. Breaking the loop stops further HTTP requests.
func (i *Iterator[T]) All() iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		for i.Next() {
			if !yield(i.Current(), nil) {
				return
			}
		}
		if err := i.Err(); err != nil {
			yield(nil, err)
		}
	}
}
