package validator

import (
	"fmt"
	"regexp"
)

// StringValidator validates string constraints (minLength, enum, regex).
type StringValidator struct{}

// Validate performs string validation.
func (sv *StringValidator) Validate(value interface{}, path string, rule Rule) error {
	// Skip if no string constraints
	if rule.MinLength == nil && len(rule.Enum) == 0 && rule.Regex == "" {
		return nil
	}
	
	str, ok := value.(string)
	if !ok {
		if rule.MinLength != nil {
			return fmt.Errorf("path '%s': minLength validation requires a string, got %T", path, value)
		}
		if len(rule.Enum) > 0 {
			return fmt.Errorf("path '%s': enum validation requires a string, got %T", path, value)
		}
		if rule.Regex != "" {
			return fmt.Errorf("path '%s': regex validation requires a string, got %T", path, value)
		}
		return nil
	}
	
	// MinLength validation
	if rule.MinLength != nil {
		if len(str) < *rule.MinLength {
			return fmt.Errorf("path '%s': length %d is less than minimum length %d", path, len(str), *rule.MinLength)
		}
	}
	
	// Enum validation
	if len(rule.Enum) > 0 {
		match := false
		for _, enumVal := range rule.Enum {
			if str == enumVal {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("path '%s': value '%s' is not in the allowed list %v", path, str, rule.Enum)
		}
	}
	
	// Regex validation
	if rule.Regex != "" {
		compiledRegex, err := regexp.Compile(rule.Regex)
		if err != nil {
			return fmt.Errorf("path '%s': invalid regex pattern '%s': %w", path, rule.Regex, err)
		}
		if !compiledRegex.MatchString(str) {
			return fmt.Errorf("path '%s': value '%s' does not match regex pattern '%s'", path, str, rule.Regex)
		}
	}
	
	return nil
}
