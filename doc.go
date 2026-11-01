// Package raims provides typed services for the Raims public API.
//
// Reuse a Client across requests and pass a context.Context to every call.
// Amounts use int64 cents and timestamps use Unix seconds. Request fields use
// Optional to distinguish omission, null and explicit values. Writes are never
// retried automatically, including when an idempotency key is supplied.
package raims
