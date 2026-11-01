package raims

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// DefaultWebhookTolerance is the allowed difference from the signing timestamp.
const DefaultWebhookTolerance = 5 * time.Minute

// SignatureError reports an invalid webhook signature, timestamp, or body.
// It does not retain the payload, signing headers, or secret.
type SignatureError struct {
	message string
}

// Error describes the signature failure without retaining the secret or payload.
func (e *SignatureError) Error() string { return e.message }

var webhookDeliveryIDPattern = regexp.MustCompile(`\A[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\z`)

// VerifyWebhookSignature authenticates the original HTTP body bytes.
// A zero now uses the current time. Zero tolerance requires the same Unix second;
// it does not disable timestamp verification. Headers must occur exactly once.
// The caller must also check the event environment and durably deduplicate its ID.
func VerifyWebhookSignature(payload []byte, headers http.Header, secret string, tolerance time.Duration, now time.Time) error {
	if tolerance < 0 {
		return &SignatureError{message: "Webhook tolerance must be non-negative."}
	}
	if now.IsZero() {
		now = time.Now()
	}
	current := now.Unix()
	if current < 0 {
		return &SignatureError{message: "The current webhook time must be non-negative."}
	}
	id, err := readWebhookHeader(headers, "webhook-id")
	if err != nil {
		return err
	}
	timestamp, err := readWebhookHeader(headers, "webhook-timestamp")
	if err != nil {
		return err
	}
	signatures, err := readWebhookHeader(headers, "webhook-signature")
	if err != nil {
		return err
	}
	if !webhookDeliveryIDPattern.MatchString(id) {
		return &SignatureError{message: "The webhook delivery ID is invalid."}
	}
	if len(timestamp) == 0 || len(timestamp) > 12 {
		return &SignatureError{message: "The webhook timestamp is invalid."}
	}
	for i := range timestamp {
		if timestamp[i] < '0' || timestamp[i] > '9' {
			return &SignatureError{message: "The webhook timestamp is invalid."}
		}
	}
	signedAt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return &SignatureError{message: "The webhook timestamp is invalid."}
	}
	difference := current - signedAt
	if signedAt > current {
		difference = signedAt - current
	}
	if difference > int64(tolerance/time.Second) {
		return &SignatureError{message: "The webhook timestamp is outside the allowed tolerance."}
	}
	key, err := decodeWebhookSecret(secret)
	if err != nil {
		return err
	}
	defer clear(key[:])
	mac := hmac.New(sha256.New, key[:])
	_, _ = mac.Write([]byte(id + "." + timestamp + "."))
	_, _ = mac.Write(payload)
	expected := mac.Sum(nil)
	matched := false
	for _, entry := range strings.Split(signatures, ",") {
		candidate := strings.Trim(entry, " \t")
		if len(candidate) != 67 || !strings.EqualFold(candidate[:3], "v1=") {
			continue
		}
		received, decodeErr := hex.DecodeString(candidate[3:])
		if decodeErr == nil && hmac.Equal(expected, received) {
			matched = true
		}
	}
	if !matched {
		return &SignatureError{message: "No webhook signature matches the configured secret."}
	}
	return nil
}

// ConstructWebhookEvent authenticates the original body before decoding its JSON
// object. Raw values preserve exact integers and unknown event fields.
func ConstructWebhookEvent(payload []byte, headers http.Header, secret string, tolerance time.Duration, now time.Time) (map[string]json.RawMessage, error) {
	if err := VerifyWebhookSignature(payload, headers, secret, tolerance, now); err != nil {
		return nil, err
	}
	if !utf8.Valid(payload) {
		return nil, &SignatureError{message: "The webhook body is not valid UTF-8."}
	}
	var event map[string]json.RawMessage
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, &SignatureError{message: "The webhook body must be a valid JSON object."}
	}
	if event == nil {
		return nil, &SignatureError{message: "The webhook body must be a JSON object."}
	}
	return event, nil
}

func readWebhookHeader(headers http.Header, name string) (string, error) {
	found := false
	value := ""
	for key, values := range headers {
		if !strings.EqualFold(key, name) {
			continue
		}
		if found {
			return "", &SignatureError{message: "A webhook signing header occurs more than once."}
		}
		found = true
		if len(values) != 1 || values[0] == "" {
			return "", &SignatureError{message: "A required webhook signing header is missing or invalid."}
		}
		value = values[0]
		for i := range value {
			if (value[i] < 32 && value[i] != '\t') || value[i] == 127 {
				return "", &SignatureError{message: "A required webhook signing header is missing or invalid."}
			}
		}
	}
	if !found {
		return "", &SignatureError{message: "A required webhook signing header is missing or invalid."}
	}
	return value, nil
}

func decodeWebhookSecret(secret string) ([32]byte, error) {
	var decoded [32]byte
	const prefix = "whsec_"
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	invalid := func() ([32]byte, error) {
		clear(decoded[:])
		return decoded, &SignatureError{message: "The webhook signing secret is invalid."}
	}
	if !strings.HasPrefix(secret, prefix) {
		return invalid()
	}
	encoded := secret[len(prefix):]
	if len(encoded) == 0 || len(encoded) > 43 || (len(encoded) > 1 && encoded[0] == '0') {
		return invalid()
	}
	for i := range encoded {
		digit := strings.IndexByte(alphabet, encoded[i])
		if digit < 0 {
			return invalid()
		}
		carry := uint32(digit)
		for offset := range decoded {
			carry += uint32(decoded[offset]) * 62
			decoded[offset] = byte(carry)
			carry >>= 8
		}
		if carry != 0 {
			return invalid()
		}
	}
	return decoded, nil
}
