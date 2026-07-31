package validator

import (
	"fmt"
	"reflect"
)

// TypeValidator validates the type of a value.
type TypeValidator struct{}

// normalizeTypeName maps common schema type names to Go reflect.Kind strings.
func normalizeTypeName(t string) string {
	switch t {
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	case "array":
		return "slice"
	case "object":
		return "map"
	case "number", "float", "double":
		return "number"
	default:
		return t
	}
}

// typeKinds maps a normalized logical type name to every Go kind that satisfies
// it. Parsers disagree on the concrete width they produce for the same literal
// (YAML gives int, JSON and TOML give int64), so a type rule has to accept the
// whole family rather than one exact kind.
var typeKinds = map[string][]reflect.Kind{
	"bool":   {reflect.Bool},
	"string": {reflect.String},
	"int": {
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	},
	"slice": {reflect.Slice, reflect.Array},
	"map":   {reflect.Map},
}

// Validate performs type validation.
func (tv *TypeValidator) Validate(value interface{}, path string, rule Rule) error {
	if rule.Type == "" {
		return nil // No type validation required
	}

	if value == nil {
		return fmt.Errorf("path '%s': expected type %s, got null", path, rule.Type)
	}

	normalizedType := normalizeTypeName(rule.Type)

	valKind := reflect.TypeOf(value).Kind()

	// Handle number type (supports all Go numeric types internally)
	if normalizedType == "number" {
		if _, ok := NumberFromInterface(value); !ok {
			return fmt.Errorf("path '%s': expected type %s, got %T", path, rule.Type, value)
		}
		return nil
	}

	// Match against the set of Go kinds that satisfy the logical type rather than
	// comparing kind names. Each parser produces its own concrete width — YAML
	// yields int, JSON and TOML yield int64 — so a name comparison made the same
	// rule pass for one input format and fail for another.
	if kinds, known := typeKinds[normalizedType]; known {
		for _, k := range kinds {
			if valKind == k {
				return nil
			}
		}
		return fmt.Errorf("path '%s': expected type %s, got %s", path, rule.Type, valKind)
	}

	// Unrecognised type names keep the original kind-name comparison so any
	// schema relying on an exact Go kind (e.g. "int64") behaves as before.
	if valKind.String() != normalizedType {
		return fmt.Errorf("path '%s': expected type %s, got %s", path, rule.Type, valKind)
	}

	return nil
}
