# Configuration and types

## Client and transport

```go
client, err := raims.NewClient(apiKey,
    raims.WithTimeout(80*time.Second),
    raims.WithMaxNetworkRetries(2),
)
```

The default timeout is 80 seconds per attempt, including response reading.
The default transport sets a ten-second connection timeout and reuses
connections. Pass a context with a deadline to limit the entire call, including
retries and waits. `WithRequestTimeout` adjusts the timeout for a call.

`WithHTTPClient` accepts an application-provided `*http.Client`. The SDK copies
its configuration, rejects redirects and does not use cookies. The supplied
transport must respect the context and preserve the policy of sending writes
only once. `CloseIdleConnections` closes idle connections in the client's
transport. If you share a transport, this call also affects its idle connections.

`WithBaseURL` accepts a trusted HTTPS origin without a path, query or
credentials. HTTP is accepted only for loopback IP addresses during local
validation. The origin receives the credential; use only a destination you control.

## API keys and OAuth

`NewClient` uses an API key. The key determines the environment; the client
does not offer environment selection or send `X-Raims-Environment` in this mode.

For an application authorized through OAuth:

```go
client, err := raims.NewOAuthClient(accessToken, raims.Sandbox,
    []raims.Environment{raims.Sandbox, raims.Production},
)
```

Provide only the environments actually granted to the connection. The SDK
accepts an access token obtained through the application's flow; the application
remains responsible for authorization, token storage and renewal. Create another
client when you renew the token or change environments.

The client sends the selected environment and rejects incompatible endpoints
before sending HTTP. Simulations accept only Sandbox. The environment never
changes automatically. The API remains responsible for checking the credential,
resource and permissions.

## Omission, null and values

Input fields use `Optional[T]`. Its zero value omits the field.
`Set(value)` sends the value, including `false`, zero or an empty string.
`Null[T]()` sends null when the field allows it.

```go
params := &raims.UpdateCustomerParams{
    FullName: raims.Null[string](),
    Metadata: raims.Set(map[string]*string{}),
}
```

This example removes the name and preserves metadata by sending `{}`. To clear
the entire map, use `raims.Null[map[string]*string]()`. To remove a key, send nil
as its value. Use an initialized empty map or slice to send `{}` or `[]`;
a nil map or slice is serialized as null.

Enum constants are generated from the contract. Integers use `int64`, including
numeric constants. For literals, use `Set[int64](15000)` to match the field
without inferring `int`.

## Responses

Responses have typed fields, `LastResponse` and `RawJSON`. `HasField`
distinguishes an absent field; `IsNull` identifies explicit null. `ExtraFields`
preserves future fields and `RawJSON` stores the original response. Use
`RawJSON` when you need to retain the exact representation received.

`json.Marshal` preserves future fields, omission and null in received responses.
It also includes changes made to typed fields and `ExtraFields`. `HasField` and
`IsNull` describe the response as received; they do not track later changes to
the struct. JSON field names are matched exactly, so an unknown field cannot
replace a known property.

The bytes in `ExtraFields` are independent of `RawJSON`. You can edit an extra
value or append bytes to it without changing the received snapshot or the
presence and null indicators.

Structural validation checks required fields, types and nullability. Amounts
in cents and dates do not pass through floating-point values. New enum values
in responses remain readable; decide how to handle them in your application.

`LastResponse` retains the status, `RequestID`, `RetryAfter` and
`IdempotencyReplayed`. A `202` status remains pending, without automatic polling.
`Charges.Revoke` uses `RevokeChargeResult` to represent both `200` and `202`.

Document and upload-preparation methods return the data provided by the API.
The SDK does not download the returned URLs or perform uploads with the API key.

## Diagnostics

`WithRequestObserver` receives the method, path without IDs, attempt, duration,
status and Request-Id after each attempt. The callback may be called by different
goroutines and must be safe for concurrent use. A callback failure does not
change the call's result.

The callback runs synchronously and should return quickly. Its execution time
does not discard a conclusive response already received, including for writes.
Context deadlines continue to control the transport and any eligible retries.

Aggregate connection, HTTP and protocol errors to detect outages or incompatible
response changes. Confirm recovery with a new valid read. Avoid logging entire
objects, payloads and temporary URLs.
