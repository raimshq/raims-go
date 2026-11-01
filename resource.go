package raims

import (
	"encoding/json"
	"fmt"
	"time"
)

// Response contains metadata only. Status 202 indicates a pending result.
type Response struct {
	// StatusCode is the received HTTP status; 202 indicates a pending result.
	StatusCode int
	// RequestID is the response correlation identifier, when available.
	RequestID string
	// RetryAfter is the parsed scheduling hint, capped at 24 hours.
	RetryAfter time.Duration
	// IdempotencyReplayed reports the API's Idempotency-Replayed header.
	IdempotencyReplayed bool
}

// Resource retains the original fields, including future fields and explicit nulls.
type Resource struct {
	// LastResponse contains HTTP metadata for a resource returned by a client call.
	LastResponse *Response `json:"-"`
	// RawJSON preserves the original bytes of this object, including future fields.
	RawJSON json.RawMessage `json:"-"`
	// ExtraFields contains editable future fields with bytes independent of RawJSON.
	ExtraFields map[string]json.RawMessage `json:"-"`
	fields      map[string]json.RawMessage
	pageRequest *pageRequest
}

func (r *Resource) setResponse(response *Response) { r.LastResponse = response }

// HasField reports whether the received object contained the exact JSON property name.
// Edits to typed properties or ExtraFields do not change this original presence metadata.
func (r Resource) HasField(name string) bool { _, ok := r.fields[name]; return ok }

// IsNull reports whether the received property was explicitly null.
// Omitted properties return false; edits to typed properties or ExtraFields do not change it.
func (r Resource) IsNull(name string) bool {
	value, ok := r.fields[name]
	return ok && string(value) == "null"
}

// APIError is a non-success HTTP response. Inspect Code and Detail for handling.
// Error intentionally excludes response payloads and messages from logs.
type APIError struct {
	// StatusCode is the failed HTTP status.
	StatusCode int
	// Code is the structured API code, or empty when no valid error envelope was received.
	Code string
	// Detail contains the structured error when the response envelope is valid.
	Detail *Error
	// Response contains the received HTTP metadata.
	Response *Response
}

// Error returns the HTTP failure status without including response payloads.
func (e *APIError) Error() string { return fmt.Sprintf("raims: API returned HTTP %d", e.StatusCode) }

// ConnectionError indicates that no conclusive response was received.
// A write may already have been processed. Unwrap supports errors.Is.
type ConnectionError struct{ cause error }

// Error describes the connection failure without including the request URL or payload.
func (e *ConnectionError) Error() string { return "raims: connection failed or timed out" }

// Unwrap exposes the cause for errors.Is and errors.As, including context deadlines.
// Inspect the underlying error before logging it because it may contain request details.
func (e *ConnectionError) Unwrap() error { return e.cause }

// ProtocolError indicates an unexpected status, body, or pagination cursor.
type ProtocolError struct {
	// Reason describes the protocol mismatch without including payloads.
	Reason string
	// Response contains HTTP metadata when a response was received.
	Response *Response
}

// Error returns the protocol mismatch description.
func (e *ProtocolError) Error() string { return "raims: " + e.Reason }

// ValidationError is detected before sending a request.
type ValidationError struct {
	// Field identifies the argument or request property to correct.
	Field string
	// Reason describes the rejected argument without including its value.
	Reason string
}

// Error identifies the invalid argument and the required correction.
func (e *ValidationError) Error() string { return "raims: " + e.Field + " " + e.Reason }
