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

// Validate performs type validation.
func (tv *TypeValidator) Validate(value interface{}, path string, rule Rule) error {
	if rule.Type == "" {
		return nil // No type validation required
	}

	if value == nil {
		return fmt.Errorf("path '%s': expected type %s, got null", path, rule.Type)
	}

	normalizedType := normalizeTypeName(rule.Type)

	valType := reflect.TypeOf(value).Kind().String()

	// Handle number type (supports all Go numeric types internally)
	if normalizedType == "number" {
		if _, ok := NumberFromInterface(value); !ok {
			return fmt.Errorf("path '%s': expected type %s, got %T", path, rule.Type, value)
		}
		return nil
	}

	if valType != normalizedType {
		return fmt.Errorf("path '%s': expected type %s, got %s", path, rule.Type, valType)
	}

	return nil
}
