package raims

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Version identifies this SDK in the HTTP User-Agent.
const Version = "0.1.2"

// DefaultBaseURL is the origin of the Raims public API.
const DefaultBaseURL = "https://api.raims.com"

// Environment selects the data environment of an OAuth client.
// API keys determine their own environment and do not use this setting.
type Environment string

// Client is reusable by concurrent goroutines. Its configuration is immutable.
type Client struct {
	// Accounts provides account context, fees and checkout settings.
	Accounts *AccountsService
	// Balances provides the account balance and forecast.
	Balances *BalancesService
	// Charges provides charges, payments, documents and charge refunds.
	Charges *ChargesService
	// Customers provides customer registration and retrieval.
	Customers *CustomersService
	// Disputes provides dispute retrieval and listing.
	Disputes *DisputesService
	// PaymentLinks provides payment links and requests to send them by email.
	PaymentLinks *PaymentLinksService
	// Refunds provides refund retrieval, listing and receipts.
	Refunds *RefundsService
	// Statements provides financial statement entries.
	Statements *StatementsService
	// WebhookEventTypes lists the event types available for webhook subscriptions.
	WebhookEventTypes *WebhookEventTypesService
	// Webhooks provides webhook subscriptions, health, pause, tests and secret rotation.
	Webhooks *WebhooksService
	// WithdrawalDestinations provides the configured withdrawal destinations.
	WithdrawalDestinations *WithdrawalDestinationsService
	// Withdrawals provides withdrawals, limits, summaries and automatic settings.
	Withdrawals *WithdrawalsService
	// Simulations provides charge and refund simulations in Sandbox.
	Simulations *SimulationsService
	key         string
	baseURL     string
	http        *http.Client
	timeout     time.Duration
	retries     int
	oauth       bool
	environment Environment
	authorized  []Environment
	observer    func(RequestDiagnostics)
}

// RequestDiagnostics contains no credentials, resource IDs, queries or payloads.
type RequestDiagnostics struct {
	// Method is the HTTP method of the attempt.
	Method string
	// Path is the endpoint template, with placeholders instead of resource IDs.
	Path string
	// Attempt starts at one and increases for each eligible read retry.
	Attempt int
	// Duration covers transport and response-body reading, before the observer runs.
	Duration time.Duration
	// StatusCode is zero when no HTTP response was received.
	StatusCode int
	// RequestID is the response correlation identifier, when available.
	RequestID string
}

type clientConfig struct {
	baseURL  string
	http     *http.Client
	timeout  time.Duration
	retries  int
	observer func(RequestDiagnostics)
}

// ClientOption configures a client through one of the With helpers.
// Options are applied and validated when the client is constructed.
type ClientOption func(*clientConfig) error

// WithBaseURL sets a trusted HTTPS origin receiving the credential.
// An HTTP loopback IP is also accepted for local validation. Paths, queries,
// fragments and credentials in the URL are rejected during construction.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) error { c.baseURL = baseURL; return nil }
}

// WithHTTPClient copies the supplied HTTP client configuration and shares its transport.
// The SDK refuses redirects and disables cookies. The transport must respect
// request contexts and preserve the single-send policy for writes.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *clientConfig) error {
		if client == nil {
			return &ValidationError{Field: "HTTPClient", Reason: "must not be nil"}
		}
		c.http = client
		return nil
	}
}

// WithTimeout sets the positive timeout of each HTTP attempt, including body reads.
// The default is 80 seconds. A context deadline can limit the entire call,
// including retries and backoff; WithRequestTimeout overrides an individual call.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) error {
		if timeout <= 0 {
			return &ValidationError{Field: "timeout", Reason: "must be positive"}
		}
		c.timeout = timeout
		return nil
	}
}

// WithMaxNetworkRetries sets zero to five additional GET attempts; the default is zero.
// Eligible failures are connection errors and HTTP 408, 429, 500, 502, 503 and 504.
// Writes are never retried, including when an idempotency key is supplied.
func WithMaxNetworkRetries(retries int) ClientOption {
	return func(c *clientConfig) error {
		if retries < 0 || retries > 5 {
			return &ValidationError{Field: "retries", Reason: "must be between zero and five"}
		}
		c.retries = retries
		return nil
	}
}

// WithRequestObserver registers a synchronous callback after each transport attempt.
// The callback must return promptly and support concurrent calls. A nil callback
// disables observation. Its panic or duration does not discard a conclusive response.
func WithRequestObserver(observer func(RequestDiagnostics)) ClientOption {
	return func(c *clientConfig) error { c.observer = observer; return nil }
}

// NewClient creates a client authenticated with an API key. The key defines its environment.
func NewClient(apiKey string, options ...ClientOption) (*Client, error) {
	return newClient(apiKey, options...)
}

// NewOAuthClient uses an access token obtained by the application's OAuth flow.
// The caller supplies the environments granted by the connection. The SDK does not renew tokens.
func NewOAuthClient(accessToken string, environment Environment, authorized []Environment, options ...ClientOption) (*Client, error) {
	if !contains([]Environment{Sandbox, Production}, environment) || !contains(authorized, environment) {
		return nil, &ValidationError{Field: "environment", Reason: "must be granted by the OAuth connection"}
	}
	for _, value := range authorized {
		if !contains([]Environment{Sandbox, Production}, value) {
			return nil, &ValidationError{Field: "authorized environments", Reason: "contains an invalid value"}
		}
	}
	client, err := newClient(accessToken, options...)
	if err != nil {
		return nil, err
	}
	client.oauth = true
	client.environment = environment
	client.authorized = append([]Environment(nil), authorized...)
	return client, nil
}

func newClient(key string, options ...ClientOption) (*Client, error) {
	if key == "" || !visibleASCII(key) {
		return nil, &ValidationError{Field: "credential", Reason: "must contain visible ASCII characters"}
	}
	config := clientConfig{baseURL: DefaultBaseURL, timeout: 80 * time.Second}
	for _, option := range options {
		if option == nil {
			return nil, &ValidationError{Field: "option", Reason: "must not be nil"}
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}
	base, err := url.Parse(config.baseURL)
	if err != nil || base.Host == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" || (base.Path != "" && base.Path != "/") || base.Opaque != "" {
		return nil, &ValidationError{Field: "base URL", Reason: "must be an HTTPS origin"}
	}
	if base.Scheme != "https" {
		ip := net.ParseIP(base.Hostname())
		if base.Scheme != "http" || ip == nil || !ip.IsLoopback() {
			return nil, &ValidationError{Field: "base URL", Reason: "requires HTTPS, except for loopback IPs"}
		}
	}
	var httpClient http.Client
	if config.http != nil {
		httpClient = *config.http
	} else {
		transport := &http.Transport{Proxy: http.ProxyFromEnvironment, ForceAttemptHTTP2: true, MaxIdleConns: 100, IdleConnTimeout: 90 * time.Second, TLSHandshakeTimeout: 10 * time.Second, ExpectContinueTimeout: time.Second}
		if standard, ok := http.DefaultTransport.(*http.Transport); ok {
			transport = standard.Clone()
		}
		transport.DialContext = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext
		httpClient.Transport = transport
	}
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	httpClient.Jar = nil
	client := &Client{key: key, baseURL: strings.TrimSuffix(base.String(), "/"), http: &httpClient, timeout: config.timeout, retries: config.retries, observer: config.observer}
	initServices(client)
	return client, nil
}

// String returns a client identifier without exposing its credential.
func (c *Client) String() string { return "raims.Client" }

// GoString returns a client identifier for Go formatting without exposing its credential.
func (c *Client) GoString() string { return "raims.Client" }

// CloseIdleConnections closes idle connections in the client's transport.
// Active requests continue. A shared transport also closes its other idle connections.
func (c *Client) CloseIdleConnections() { c.http.CloseIdleConnections() }

type requestConfig struct {
	idempotencyKey string
	timeout        time.Duration
}

// RequestOption configures a single request and is validated before HTTP is sent.
// Collection navigation preserves the options of its original request.
type RequestOption func(*requestConfig) error

// WithIdempotencyKey sends a persisted key for an endpoint that accepts idempotency.
// The key must contain 1 to 255 visible ASCII characters without spaces.
// Keep the same key and parameters when repeating the same intention.
func WithIdempotencyKey(key string) RequestOption {
	return func(c *requestConfig) error {
		if len(key) < 1 || len(key) > 255 || !visibleASCII(key) {
			return &ValidationError{Field: "Idempotency-Key", Reason: "must contain 1 to 255 visible ASCII characters"}
		}
		c.idempotencyKey = key
		return nil
	}
}

// WithRequestTimeout overrides the positive timeout of each HTTP attempt for this call.
// It includes response-body reading. Use a context deadline to bound retries and backoff.
func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(c *requestConfig) error {
		if timeout <= 0 {
			return &ValidationError{Field: "timeout", Reason: "must be positive"}
		}
		c.timeout = timeout
		return nil
	}
}
func visibleASCII(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] < 33 || value[i] > 126 {
			return false
		}
	}
	return true
}

func request[T any](ctx context.Context, client *Client, key string, path map[string]string, params any, options ...RequestOption) (*T, error) {
	query, err := marshalParams(params)
	if err != nil {
		return nil, err
	}
	if operations[key].BodyRequired {
		raw, err := json.Marshal(params)
		if err != nil || bytes.Equal(raw, []byte("null")) {
			return nil, &ValidationError{Field: "params", Reason: "is required"}
		}
	}
	var result T
	metadata, err := client.execute(ctx, key, path, query, &result, options...)
	if err != nil {
		return nil, err
	}
	if resource, ok := any(&result).(interface{ setResponse(*Response) }); ok {
		resource.setResponse(metadata)
	}
	if resource, ok := any(&result).(interface {
		pagination() pageInfo
		setPageRequest(*pageRequest)
	}); ok {
		resource.setPageRequest(&pageRequest{client: client, key: key, path: clonePath(path), params: query, options: append([]RequestOption(nil), options...)})
	}
	return &result, nil
}

func marshalParams(params any) (map[string]json.RawMessage, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, &ValidationError{Field: "params", Reason: "cannot be serialized"}
	}
	if bytes.Equal(data, []byte("null")) {
		return map[string]json.RawMessage{}, nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, &ValidationError{Field: "params", Reason: "must be an object"}
	}
	return fields, nil
}

func (c *Client) execute(ctx context.Context, key string, paths map[string]string, params any, result any, options ...RequestOption) (*Response, error) {
	if ctx == nil {
		return nil, &ValidationError{Field: "context", Reason: "must not be nil"}
	}
	op, ok := operations[key]
	if !ok {
		return nil, &ValidationError{Field: "endpoint", Reason: "is not in the pinned contract"}
	}
	if c.oauth && (!contains(op.Environments, c.environment) || !contains(c.authorized, c.environment)) {
		return nil, &ValidationError{Field: "environment", Reason: "is not accepted by this endpoint and OAuth connection"}
	}
	config := requestConfig{timeout: c.timeout}
	for _, option := range options {
		if option == nil {
			return nil, &ValidationError{Field: "request option", Reason: "must not be nil"}
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}
	if config.idempotencyKey != "" && !op.Idempotency {
		return nil, &ValidationError{Field: "Idempotency-Key", Reason: "is not accepted by this endpoint"}
	}
	path := op.Path
	for name, spec := range op.PathParameters {
		value := paths[name]
		if value == "" || !visibleASCII(value) || strings.ContainsAny(value, "/\\?#%") || value == "." || value == ".." {
			return nil, &ValidationError{Field: "path parameter", Reason: "is invalid"}
		}
		if spec.Pattern != "" {
			expression, err := regexp.Compile(spec.Pattern)
			if err != nil || !expression.MatchString(value) {
				return nil, &ValidationError{Field: "path parameter", Reason: "does not match the public format"}
			}
		}
		path = strings.ReplaceAll(path, "{"+name+"}", url.PathEscape(value))
	}
	fields, err := marshalParams(params)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(fields)
	if err != nil {
		return nil, &ValidationError{Field: "params", Reason: "cannot be serialized"}
	}
	if op.Params != "" {
		value, err := parseJSON(data)
		if err != nil {
			return nil, &ValidationError{Field: "params", Reason: "is not valid JSON"}
		}
		if err := validateValue(value, schemas[op.Params], true, "params"); err != nil {
			return nil, err
		}
	}
	if op.BodyRequired && params == nil {
		return nil, &ValidationError{Field: "params", Reason: "is required"}
	}
	endpoint := c.baseURL + path
	var body []byte
	if op.Body {
		body = data
		if op.MaxBodyBytes > 0 && len(body) > op.MaxBodyBytes {
			return nil, &ValidationError{Field: "body", Reason: "exceeds the endpoint limit"}
		}
	} else {
		query := url.Values{}
		for name, raw := range fields {
			var value any
			decoder := json.NewDecoder(bytes.NewReader(raw))
			decoder.UseNumber()
			if err := decoder.Decode(&value); err != nil {
				return nil, &ValidationError{Field: "query", Reason: "cannot be serialized"}
			}
			query.Set(name, fmt.Sprint(value))
		}
		if len(query) > 0 {
			endpoint += "?" + query.Encode()
		}
	}
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		started := time.Now()
		attemptCtx, cancel := context.WithTimeout(ctx, config.timeout)
		var requestBody io.Reader = bytes.NewReader(body)
		if op.Method != http.MethodGet {
			// A non-replayable body also prevents transport-level retries of writes
			// when Idempotency-Key is present, including writes with an empty body.
			requestBody = io.NopCloser(bytes.NewReader(body))
		}
		req, err := http.NewRequestWithContext(attemptCtx, op.Method, endpoint, requestBody)
		if err != nil {
			cancel()
			return nil, &ValidationError{Field: "request", Reason: "cannot be constructed"}
		}
		if op.Method != http.MethodGet {
			req.GetBody = nil
			req.ContentLength = int64(len(body))
			if len(body) == 0 {
				req.ContentLength = -1
			}
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "RaimsGo/"+Version)
		if op.Body {
			req.Header.Set("Content-Type", "application/json")
		}
		if config.idempotencyKey != "" {
			req.Header.Set("Idempotency-Key", config.idempotencyKey)
		}
		if c.oauth {
			req.Header.Set("X-Raims-Environment", string(c.environment))
		}
		response, callErr := c.http.Do(req)
		var responseData []byte
		var metadata *Response
		if response != nil {
			metadata = responseMetadata(response)
			if callErr == nil {
				responseData, callErr = io.ReadAll(io.LimitReader(response.Body, (16<<20)+1))
				if len(responseData) > 16<<20 {
					callErr = &ProtocolError{Reason: "response exceeds 16 MiB", Response: metadata}
				}
			}
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}
		cancel()
		diagnostic := RequestDiagnostics{Method: op.Method, Path: op.Path, Attempt: attempt + 1, Duration: time.Since(started)}
		if metadata != nil {
			diagnostic.StatusCode = metadata.StatusCode
			diagnostic.RequestID = metadata.RequestID
		}
		// The observer runs after the transport attempt. Its duration must not
		// change the cancellation decision for the response already received.
		contextErr := ctx.Err()
		c.observe(diagnostic)
		if contextErr != nil {
			return nil, contextErr
		}
		if callErr != nil {
			var protocol *ProtocolError
			if errors.As(callErr, &protocol) {
				return nil, callErr
			}
		}
		retry := op.Method == http.MethodGet && attempt < c.retries && (callErr != nil || (metadata != nil && contains([]int{408, 429, 500, 502, 503, 504}, metadata.StatusCode)))
		if retry && (metadata == nil || metadata.RetryAfter <= 30*time.Second) {
			delay := retryDelay(attempt)
			if metadata != nil && metadata.RetryAfter > delay {
				delay = metadata.RetryAfter
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
			continue
		}
		if callErr != nil {
			return nil, &ConnectionError{cause: callErr}
		}
		if metadata == nil {
			return nil, &ProtocolError{Reason: "transport returned no response"}
		}
		if metadata.StatusCode >= 300 && metadata.StatusCode < 400 {
			return nil, &ProtocolError{Reason: "redirect was refused", Response: metadata}
		}
		if metadata.StatusCode < 200 || metadata.StatusCode >= 300 {
			failure := &APIError{StatusCode: metadata.StatusCode, Response: metadata}
			var envelope ErrorResponse
			if json.Unmarshal(responseData, &envelope) == nil {
				failure.Detail = &envelope.Error
				failure.Code = envelope.Error.Code
			}
			return nil, failure
		}
		schemaName, expected := op.Responses[strconv.Itoa(metadata.StatusCode)]
		if !expected {
			return nil, &ProtocolError{Reason: "unexpected success status", Response: metadata}
		}
		if schemaName == "" {
			if len(bytes.TrimSpace(responseData)) != 0 {
				return nil, &ProtocolError{Reason: "unexpected response body", Response: metadata}
			}
			return metadata, nil
		}
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(response.Header.Get("Content-Type"), ";")[0]))
		if mediaType != "application/json" && !strings.HasSuffix(mediaType, "+json") {
			return nil, &ProtocolError{Reason: "success response is not JSON", Response: metadata}
		}
		if result != nil {
			if _, err := decodeResource(responseData, schemaName, result); err != nil {
				return nil, protocolResponse(err, metadata)
			}
		}
		return metadata, nil
	}
}

func responseMetadata(response *http.Response) *Response {
	requestID := response.Header.Get("Request-Id")
	return &Response{StatusCode: response.StatusCode, RequestID: requestID, RetryAfter: parseRetryAfter(response.Header.Get("Retry-After"), time.Now()), IdempotencyReplayed: strings.EqualFold(response.Header.Get("Idempotency-Replayed"), "true")}
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	value = strings.TrimSpace(value)
	if value != "" && strings.Trim(value, "0123456789") == "" {
		trimmed := strings.TrimLeft(value, "0")
		if len(trimmed) > 5 {
			return 24 * time.Hour
		}
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds >= 0 {
		if seconds > 86400 {
			return 24 * time.Hour
		}
		return time.Duration(seconds) * time.Second
	}
	if date, err := http.ParseTime(value); err == nil && date.After(now) {
		if date.Sub(now) > 24*time.Hour {
			return 24 * time.Hour
		}
		return date.Sub(now)
	}
	return 0
}

func protocolResponse(err error, metadata *Response) error {
	var protocol *ProtocolError
	if errors.As(err, &protocol) {
		protocol.Response = metadata
		return protocol
	}
	return &ProtocolError{Reason: "invalid response body", Response: metadata}
}

func retryDelay(attempt int) time.Duration {
	maximum := int64(250*time.Millisecond) << attempt
	value, err := rand.Int(rand.Reader, big.NewInt(maximum/2+1))
	if err != nil {
		return time.Duration(maximum)
	}
	return time.Duration(maximum/2 + value.Int64())
}

func (c *Client) observe(diagnostic RequestDiagnostics) {
	if c.observer == nil {
		return
	}
	defer func() { _ = recover() }()
	c.observer(diagnostic)
}
