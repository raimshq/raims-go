# Errors, idempotency and retries

## Errors

Use `errors.As` to check the category and `errors.Is` for cancellation:

```go
var apiError *raims.APIError
if errors.As(err, &apiError) {
    status := apiError.StatusCode
    code := apiError.Code
    requestID := apiError.Response.RequestID
    // Choose an action based on the status and code.
}
```

| Type | Meaning |
| --- | --- |
| `ValidationError` | Input rejected before HTTP |
| `APIError` | Failed HTTP response, including 429 and 5xx |
| `ConnectionError` | Interrupted response, connection failure or timeout |
| `ProtocolError` | Redirect, unexpected status, incompatible response or invalid cursor |
| `SignatureError` | Webhook delivery authentication failure |

`APIError.Detail` contains the code, message and validation details when the API
returns a valid envelope. Proxies may respond without that envelope.
`ProtocolError.Response` preserves received HTTP metadata.
The `Error()` text for HTTP and connection failures does not include the payload.
Do not log an unwrapped error, entire structs or API messages without checking
the data they contain. A connection lost during a write leaves its result
uncertain; the operation may have been processed.

## Idempotency

Pass `WithIdempotencyKey(persistedKey)` only to endpoints that accept the header,
as identified in the [reference](reference.md). The key must contain between
1 and 255 visible ASCII characters, without spaces.

Create and persist the key before sending the request, together with the
parameters for that intent. If the response is lost, retain both. The SDK does
not generate keys or retry writes. `IDEMPOTENCY_REQUEST_IN_PROGRESS` requires
waiting according to `Retry-After`; `IDEMPOTENCY_RESULT_UNKNOWN` requires lookup
and reconciliation before retrying. See the
[API rules](https://docs.raims.com/en/api/idempotencia).

## Retries

The default is zero retries. `WithMaxNetworkRetries` accepts zero to five
additional attempts, exclusively for GET. Eligible failures are connection
errors and HTTP 408, 429, 500, 502, 503 and 504.

Backoff includes jitter and honors `Retry-After`, expressed in seconds or as an
HTTP date. A wait longer than 30 seconds returns the failure so the application
can reschedule it. Each attempt has its own timeout; a context with a deadline
limits the total duration.

POST, PATCH and DELETE are sent only once, including when an idempotency key is
provided. The default transport also prevents replay of these requests when a
reused connection is lost. A custom transport must preserve that guarantee.
