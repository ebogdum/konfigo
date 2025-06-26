package validator

import (
	"fmt"
)

// NumericValidator validates numeric constraints (min/max).
type NumericValidator struct{}

// Validate performs numeric validation.
func (nv *NumericValidator) Validate(value interface{}, path string, rule Rule) error {
	// Skip if no numeric constraints
	if rule.Min == nil && rule.Max == nil {
		return nil
	}
	
	numVal, ok := NumberFromInterface(value)
	if !ok {
		return fmt.Errorf("path '%s': min/max validation requires a number, got %T", path, value)
	}
	
	num := numVal.ToFloat64()
	
	if rule.Min != nil && num < *rule.Min {
		return fmt.Errorf("path '%s': value %v is less than minimum %v", path, num, *rule.Min)
	}
	
	if rule.Max != nil && num > *rule.Max {
		return fmt.Errorf("path '%s': value %v is greater than maximum %v", path, num, *rule.Max)
	}
	
	return nil
}
