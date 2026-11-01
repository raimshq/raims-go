package raims

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

//go:embed contracts/schemas.json
var schemaBytes []byte

//go:embed contracts/operations.json
var operationBytes []byte

type schema struct {
	Ref        string             `json:"$ref"`
	Type       json.RawMessage    `json:"type"`
	Nullable   bool               `json:"nullable"`
	Properties map[string]*schema `json:"properties"`
	Required   []string           `json:"required"`
	Items      *schema            `json:"items"`
	Additional json.RawMessage    `json:"additionalProperties"`
	AnyOf      []*schema          `json:"anyOf"`
	Pattern    string             `json:"pattern"`
	kind       string
	allowsNull bool
	additional *schema
	closed     bool
}
type operation struct {
	Method         string             `json:"method"`
	Path           string             `json:"path"`
	Params         string             `json:"params"`
	Body           bool               `json:"body"`
	BodyRequired   bool               `json:"body_required"`
	Responses      map[string]string  `json:"responses"`
	Environments   []Environment      `json:"environments"`
	Idempotency    bool               `json:"idempotency"`
	PathParameters map[string]*schema `json:"path_parameters"`
	MaxBodyBytes   int                `json:"max_body_bytes"`
}

var schemas map[string]*schema
var operations map[string]operation

func init() {
	if json.Unmarshal(schemaBytes, &schemas) != nil || json.Unmarshal(operationBytes, &operations) != nil {
		panic("raims: invalid generated contract")
	}
	for _, s := range schemas {
		prepareSchema(s)
	}
}

// Compile structural metadata once; validation never reparses schema JSON.
func prepareSchema(s *schema) {
	if s == nil {
		return
	}
	var types []string
	if len(s.Type) > 0 && s.Type[0] == '[' {
		if json.Unmarshal(s.Type, &types) != nil {
			panic("raims: invalid generated type")
		}
	} else if len(s.Type) > 0 {
		var kind string
		if json.Unmarshal(s.Type, &kind) != nil {
			panic("raims: invalid generated type")
		}
		types = []string{kind}
	}
	s.allowsNull = s.Nullable
	for _, kind := range types {
		if kind == "null" {
			s.allowsNull = true
		} else {
			if s.kind != "" {
				panic("raims: unsupported generated union")
			}
			s.kind = kind
		}
	}
	s.closed = string(s.Additional) == "false"
	if len(s.Additional) > 0 && s.Additional[0] == '{' {
		if json.Unmarshal(s.Additional, &s.additional) != nil {
			panic("raims: invalid generated map")
		}
		prepareSchema(s.additional)
	}
	for _, child := range s.Properties {
		prepareSchema(child)
	}
	prepareSchema(s.Items)
	for _, branch := range s.AnyOf {
		prepareSchema(branch)
	}
}

type valueSpan struct{ start, end int }
type jsonObject struct {
	values     map[string]any
	spans      map[string]valueSpan
	source     []byte
	start, end int
}

func (o *jsonObject) raw(name string) []byte {
	span := o.spans[name]
	return o.source[span.start:span.end]
}
func parseJSON(data []byte) (any, error) {
	if !utf8.Valid(data) {
		return nil, &ProtocolError{Reason: "response is not UTF-8"}
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := readJSONValue(decoder, data, 0)
	if err != nil {
		return nil, &ProtocolError{Reason: "invalid JSON response"}
	}
	if _, err = decoder.Token(); err != io.EOF {
		return nil, &ProtocolError{Reason: "trailing JSON response data"}
	}
	return value, nil
}
func valueStart(data []byte, offset int64) int {
	start := int(offset)
	for start < len(data) {
		switch data[start] {
		case ' ', '\t', '\r', '\n', ':', ',':
			start++
		default:
			return start
		}
	}
	return start
}

// Keep exact property names and spans while rejecting duplicate keys in one pass.
func readJSONValue(decoder *json.Decoder, data []byte, depth int) (any, error) {
	if depth > 128 {
		return nil, io.ErrUnexpectedEOF
	}
	start := valueStart(data, decoder.InputOffset())
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delim {
	case '{':
		object := &jsonObject{values: map[string]any{}, spans: map[string]valueSpan{}, source: data, start: start}
		for decoder.More() {
			key, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			field, ok := key.(string)
			if !ok {
				return nil, io.ErrUnexpectedEOF
			}
			if _, exists := object.values[field]; exists {
				return nil, io.ErrUnexpectedEOF
			}
			childStart := valueStart(data, decoder.InputOffset())
			child, err := readJSONValue(decoder, data, depth+1)
			if err != nil {
				return nil, err
			}
			object.values[field] = child
			object.spans[field] = valueSpan{childStart, int(decoder.InputOffset())}
		}
		_, err := decoder.Token()
		object.end = int(decoder.InputOffset())
		return object, err
	case '[':
		array := []any{}
		for decoder.More() {
			child, err := readJSONValue(decoder, data, depth+1)
			if err != nil {
				return nil, err
			}
			array = append(array, child)
		}
		_, err := decoder.Token()
		return array, err
	default:
		return nil, io.ErrUnexpectedEOF
	}
}
func resolveSchema(s *schema) *schema {
	if s != nil && s.Ref != "" {
		return schemas[strings.TrimPrefix(s.Ref, "#/components/schemas/")]
	}
	return s
}
func concreteSchema(s *schema) *schema {
	s = resolveSchema(s)
	if s != nil && len(s.AnyOf) > 0 {
		for _, branch := range s.AnyOf {
			candidate := resolveSchema(branch)
			if candidate != nil && (candidate.kind != "" || candidate.Ref != "") {
				return concreteSchema(candidate)
			}
		}
	}
	return s
}
func acceptsNull(s *schema) bool {
	s = resolveSchema(s)
	if s == nil {
		return false
	}
	if s.allowsNull {
		return true
	}
	for _, branch := range s.AnyOf {
		if acceptsNull(branch) {
			return true
		}
	}
	return false
}
func exactInteger(number json.Number) (int64, bool) {
	text := number.String()
	if value, err := strconv.ParseInt(text, 10, 64); err == nil {
		return value, true
	}
	if len(text) > 128 {
		return 0, false
	}
	if position := strings.IndexAny(text, "eE"); position >= 0 {
		exponent, err := strconv.ParseInt(text[position+1:], 10, 32)
		if err != nil || exponent > 32 || exponent < -128 {
			return 0, false
		}
	}
	rational, ok := new(big.Rat).SetString(text)
	if !ok || !rational.IsInt() || !rational.Num().IsInt64() {
		return 0, false
	}
	return rational.Num().Int64(), true
}
func validateValue(value any, s *schema, input bool, field string) error {
	if value == nil {
		if acceptsNull(s) {
			return nil
		}
		return invalidField(input, field)
	}
	s = concreteSchema(s)
	if s == nil {
		return &ProtocolError{Reason: "unknown generated schema"}
	}
	switch s.kind {
	case "object":
		object, ok := value.(*jsonObject)
		if !ok {
			return invalidField(input, field)
		}
		for _, required := range s.Required {
			if _, ok := object.values[required]; !ok {
				return invalidField(input, field+"."+required)
			}
		}
		for key, child := range object.values {
			if property, ok := s.Properties[key]; ok {
				if err := validateValue(child, property, input, field+"."+key); err != nil {
					return err
				}
			} else if s.additional != nil {
				if err := validateValue(child, s.additional, input, field+".*"); err != nil {
					return err
				}
			} else if input && s.closed {
				return invalidField(input, field)
			}
		}
	case "array":
		array, ok := value.([]any)
		if !ok {
			return invalidField(input, field)
		}
		for _, child := range array {
			if err := validateValue(child, s.Items, input, field+"[]"); err != nil {
				return err
			}
		}
	case "string":
		if _, ok := value.(string); !ok {
			return invalidField(input, field)
		}
	case "integer":
		number, ok := value.(json.Number)
		if !ok {
			return invalidField(input, field)
		}
		if _, ok := exactInteger(number); !ok {
			return invalidField(input, field)
		}
	case "number":
		if _, ok := value.(json.Number); !ok {
			return invalidField(input, field)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return invalidField(input, field)
		}
	case "":
		return nil
	default:
		return invalidField(input, field)
	}
	return nil
}
func invalidField(input bool, field string) error {
	if input {
		return &ValidationError{Field: field, Reason: "is missing or has an incompatible type"}
	}
	return &ProtocolError{Reason: "response does not match the public contract"}
}
func contains[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

type objectBinding struct {
	fields   map[string]int
	resource int
}

var objectBindings sync.Map
var resourceType = reflect.TypeFor[Resource]()
var rawMessageType = reflect.TypeFor[json.RawMessage]()
var jsonNumberType = reflect.TypeFor[json.Number]()

func bindingFor(typ reflect.Type) *objectBinding {
	if cached, ok := objectBindings.Load(typ); ok {
		return cached.(*objectBinding)
	}
	binding := &objectBinding{fields: map[string]int{}, resource: -1}
	for index := 0; index < typ.NumField(); index++ {
		field := typ.Field(index)
		if field.Type == resourceType {
			binding.resource = index
			continue
		}
		if !field.IsExported() {
			continue
		}
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name != "" && name != "-" {
			binding.fields[name] = index
		}
	}
	actual, _ := objectBindings.LoadOrStore(typ, binding)
	return actual.(*objectBinding)
}
func resourceFor(object *jsonObject, s *schema) Resource {
	// Child resources own their bytes, so retaining one item does not retain a page.
	raw := append(json.RawMessage(nil), object.source[object.start:object.end]...)
	fields := make(map[string]json.RawMessage, len(object.spans))
	extras := map[string]json.RawMessage{}
	for name, span := range object.spans {
		end := span.end - object.start
		fields[name] = raw[span.start-object.start : end : end]
		if _, known := s.Properties[name]; !known {
			// ExtraFields is editable without changing the original response or
			// the immutable presence/null metadata stored in fields.
			extras[name] = append(json.RawMessage(nil), fields[name]...)
		}
	}
	return Resource{RawJSON: raw, ExtraFields: extras, fields: fields}
}

// Hydrate recognized JSON names directly. Unknown properties never participate
// in Go's case-insensitive struct matching, and nested models are not reparsed.
func hydrateValue(value any, s *schema, target reflect.Value, raw []byte) error {
	if !target.CanSet() {
		return &ProtocolError{Reason: "resource cannot be assigned"}
	}
	if value == nil {
		target.SetZero()
		return nil
	}
	if target.Type() == rawMessageType {
		target.SetBytes(append([]byte(nil), raw...))
		return nil
	}
	if target.Kind() == reflect.Pointer {
		target.Set(reflect.New(target.Type().Elem()))
		return hydrateValue(value, s, target.Elem(), raw)
	}
	s = concreteSchema(s)
	switch target.Kind() {
	case reflect.Struct:
		object, ok := value.(*jsonObject)
		if !ok {
			return &ProtocolError{Reason: "invalid resource object"}
		}
		binding := bindingFor(target.Type())
		for name, property := range s.Properties {
			child, present := object.values[name]
			if !present {
				continue
			}
			index, exists := binding.fields[name]
			if !exists {
				return &ProtocolError{Reason: "generated model is missing a field"}
			}
			if err := hydrateValue(child, property, target.Field(index), object.raw(name)); err != nil {
				return err
			}
		}
		if binding.resource >= 0 {
			target.Field(binding.resource).Set(reflect.ValueOf(resourceFor(object, s)))
		}
	case reflect.Map:
		object, ok := value.(*jsonObject)
		if !ok {
			return &ProtocolError{Reason: "invalid resource map"}
		}
		target.Set(reflect.MakeMapWithSize(target.Type(), len(object.values)))
		for name, child := range object.values {
			entry := reflect.New(target.Type().Elem()).Elem()
			if err := hydrateValue(child, s.additional, entry, object.raw(name)); err != nil {
				return err
			}
			target.SetMapIndex(reflect.ValueOf(name), entry)
		}
	case reflect.Slice:
		array, ok := value.([]any)
		if !ok {
			return &ProtocolError{Reason: "invalid resource array"}
		}
		target.Set(reflect.MakeSlice(target.Type(), len(array), len(array)))
		for index, child := range array {
			if err := hydrateValue(child, s.Items, target.Index(index), nil); err != nil {
				return err
			}
		}
	case reflect.String:
		if target.Type() == jsonNumberType {
			target.SetString(value.(json.Number).String())
		} else {
			target.SetString(value.(string))
		}
	case reflect.Int64:
		integer, ok := exactInteger(value.(json.Number))
		if !ok {
			return &ProtocolError{Reason: "invalid resource integer"}
		}
		target.SetInt(integer)
	case reflect.Bool:
		target.SetBool(value.(bool))
	default:
		return &ProtocolError{Reason: "unsupported generated resource type"}
	}
	return nil
}
func decodeResource(data []byte, schemaName string, destination any) (Resource, error) {
	value, err := parseJSON(data)
	if err != nil {
		return Resource{}, err
	}
	s := schemas[schemaName]
	if err := validateValue(value, s, false, "$"); err != nil {
		return Resource{}, err
	}
	target := reflect.ValueOf(destination)
	if target.Kind() != reflect.Pointer || target.IsNil() || target.Elem().Kind() != reflect.Struct {
		return Resource{}, &ProtocolError{Reason: "invalid resource destination"}
	}
	decoded := reflect.New(target.Elem().Type()).Elem()
	if err := hydrateValue(value, s, decoded, data); err != nil {
		return Resource{}, err
	}
	target.Elem().Set(decoded)
	binding := bindingFor(decoded.Type())
	if binding.resource >= 0 {
		return decoded.Field(binding.resource).Interface().(Resource), nil
	}
	return Resource{}, nil
}
func encodeResource(resource Resource, value any, schemaName string) ([]byte, error) {
	s := schemas[schemaName]
	if s == nil {
		return nil, &ProtocolError{Reason: "unknown generated schema"}
	}
	target := reflect.ValueOf(value)
	binding := bindingFor(target.Type())
	fields := make(map[string]any, len(s.Properties)+len(resource.ExtraFields))
	for name, property := range s.Properties {
		index, ok := binding.fields[name]
		if !ok {
			return nil, &ProtocolError{Reason: "generated model is missing a field"}
		}
		field := target.Field(index)
		if field.Kind() == reflect.Pointer && field.IsNil() && !contains(s.Required, name) {
			if !acceptsNull(property) || !resource.HasField(name) {
				continue
			}
		}
		fields[name] = field.Interface()
	}
	for name, raw := range resource.ExtraFields {
		if _, known := s.Properties[name]; !known {
			fields[name] = raw
		}
	}
	return json.Marshal(fields)
}
