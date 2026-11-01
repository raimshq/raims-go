package raims

import "encoding/json"

// Optional distinguishes an omitted request field from a value or explicit null.
// Use Set(value) or Null[T](). The zero value omits the field.
type Optional[T any] struct {
	value T
	set   bool
	null  bool
}

// Set includes a field with the supplied value, including false, zero and empty text.
// A nil map or slice is encoded as null; initialize an empty map or slice to send {} or [].
func Set[T any](value T) Optional[T] { return Optional[T]{value: value, set: true} }

// Null includes a field as explicit JSON null when the endpoint permits it.
func Null[T any]() Optional[T] { return Optional[T]{set: true, null: true} }

// IsZero reports whether the field is omitted. It supports the JSON omitzero tag.
func (o Optional[T]) IsZero() bool { return !o.set }

// IsSet reports whether a value or explicit null was supplied.
func (o Optional[T]) IsSet() bool { return o.set }

// IsNull reports whether the field was explicitly set to JSON null.
func (o Optional[T]) IsNull() bool { return o.set && o.null }

// Value returns the stored value and true when a non-null value was supplied.
// The boolean is false for omitted fields and explicit null.
func (o Optional[T]) Value() (T, bool) { return o.value, o.set && !o.null }

// MarshalJSON encodes the value or null. Omission is controlled by the containing
// field's omitzero tag; serializing an omitted Optional on its own produces null.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.set || o.null {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON records the supplied JSON value or explicit null.
// Invalid input leaves the existing Optional unchanged.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	var value T
	if string(data) == "null" {
		*o = Null[T]()
		return nil
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*o = Set(value)
	return nil
}
