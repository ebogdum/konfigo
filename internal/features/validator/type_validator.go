package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// TypeValidator validates the type of a value.
type TypeValidator struct{}

// Validate performs type validation.
func (tv *TypeValidator) Validate(value interface{}, path string, rule Rule) error {
	if rule.Type == "" {
		return nil // No type validation required
	}
	
	valType := reflect.TypeOf(value).Kind().String()
	
	// Handle number type (supports all Go numeric types internally)
	if rule.Type == "number" {
		if _, ok := NumberFromInterface(value); !ok {
			return fmt.Errorf("path '%s': expected type number, got %T", path, value)
		}
		return nil
	}
	
	if !strings.HasPrefix(valType, rule.Type) {
		return fmt.Errorf("path '%s': expected type %s, got %s", path, rule.Type, valType)
	}
	
	return nil
}
