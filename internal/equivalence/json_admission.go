package equivalence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type jsonAdmissionFrame struct {
	kind        json.Delim
	valueType   reflect.Type
	pendingType reflect.Type
	wantKey     bool
	names       map[string]bool
	fields      map[int]bool
}

// validateComparisonJSON runs before a decoder can discard duplicate names.
// The stack is iterative; opaque map keys are never folded. It does not add a
// whole-document allocation limit to the existing io.ReadAll loader contract.
func validateComparisonJSON(content []byte, rootType reflect.Type) error {
	if !json.Valid(content) {
		return fmt.Errorf("invalid comparison JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber() // Scanning unknown numbers must not introduce float64 limits.
	stack := []jsonAdmissionFrame{}
	for {
		token, err := decoder.Token()
		if err == io.EOF && len(stack) == 0 {
			return nil
		}
		if err != nil {
			return fmt.Errorf("invalid comparison JSON token stream")
		}
		if len(stack) > 0 {
			frame := &stack[len(stack)-1]
			if frame.kind == '{' && frame.wantKey {
				if token == json.Delim('}') {
					stack = stack[:len(stack)-1]
					continue
				}
				name, ok := token.(string)
				if !ok {
					return fmt.Errorf("invalid comparison JSON object")
				}
				if frame.names[name] {
					return fmt.Errorf("duplicate object member in comparison JSON")
				}
				frame.names[name] = true
				index, fieldType, err := comparisonJSONField(frame.valueType, name)
				if err != nil {
					return err
				}
				if index >= 0 && frame.fields[index] {
					return fmt.Errorf("ambiguous typed field in comparison JSON")
				}
				if index >= 0 {
					frame.fields[index] = true
				}
				frame.pendingType = fieldType
				frame.wantKey = false
				continue
			}
		}
		if token == json.Delim(']') {
			if len(stack) == 0 || stack[len(stack)-1].kind != '[' {
				return fmt.Errorf("invalid comparison JSON array")
			}
			stack = stack[:len(stack)-1]
			continue
		}

		valueType := rootType
		if len(stack) > 0 {
			frame := &stack[len(stack)-1]
			if frame.kind == '{' {
				valueType = frame.pendingType
				frame.wantKey = true
			} else {
				valueType = comparisonJSONArrayElement(frame.valueType)
			}
		}
		if delim, ok := token.(json.Delim); ok {
			if delim != '{' && delim != '[' {
				return fmt.Errorf("invalid comparison JSON container")
			}
			stack = append(stack, jsonAdmissionFrame{
				kind: delim, valueType: valueType, wantKey: delim == '{',
				names: map[string]bool{}, fields: map[int]bool{},
			})
		}
	}
}

func comparisonJSONIndirect(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func comparisonJSONArrayElement(t reflect.Type) reflect.Type {
	t = comparisonJSONIndirect(t)
	if t != nil && (t.Kind() == reflect.Array || t.Kind() == reflect.Slice) {
		return t.Elem()
	}
	return nil
}

// Current target DTOs have named exported fields, not promoted anonymous
// fields. Fail closed if that shape changes without a deliberate contract.
// strings.EqualFold matches pinned encoding/json v1's Unicode fold relation.
func comparisonJSONField(t reflect.Type, name string) (int, reflect.Type, error) {
	t = comparisonJSONIndirect(t)
	if t == nil {
		return -1, nil, nil
	}
	if t.Kind() == reflect.Map {
		return -1, t.Elem(), nil
	}
	if t.Kind() != reflect.Struct {
		return -1, nil, nil
	}
	folded := -1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if key == "-" {
			continue
		}
		if field.Anonymous && key == "" {
			return -1, nil, fmt.Errorf("unsupported comparison JSON target field shape")
		}
		if field.PkgPath != "" {
			continue
		}
		if key == "" {
			key = field.Name
		}
		if key == name {
			return i, field.Type, nil
		}
		if folded < 0 && strings.EqualFold(key, name) {
			folded = i
		}
	}
	if folded >= 0 {
		return folded, t.Field(folded).Type, nil
	}
	return -1, nil, nil
}
