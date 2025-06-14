package schema

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"reflect"
	"regexp"
	"strings"
)

func ApplyValidations(config map[string]interface{}, groups []ValidationGroup) error {
	logger.Debug("Applying %d validation groups...", len(groups))
	for _, group := range groups {
		val, found := util.GetNestedValue(config, group.Path)

		if group.Rules.Required && !found {
			return fmt.Errorf("path '%s' is required but not found", group.Path)
		}
		if !found {
			continue // If not required and not found, skip other rules
		}

		logger.Debug("  - Validating path '%s'", group.Path)
		rules := group.Rules

		// Type validation
		if rules.Type != "" {
			valType := reflect.TypeOf(val).Kind().String()
			// JSON numbers are float64, so we need to handle integer checks specifically
			if rules.Type == "integer" {
				if f, ok := val.(float64); !ok || f != float64(int64(f)) {
					return fmt.Errorf("path '%s': expected type integer, got %T", group.Path, val)
				}
			} else if !strings.HasPrefix(valType, rules.Type) {
				return fmt.Errorf("path '%s': expected type %s, got %s", group.Path, rules.Type, valType)
			}
		}

		// Min/Max for numbers
		if rules.Min != nil || rules.Max != nil {
			num, ok := val.(float64)
			if !ok {
				return fmt.Errorf("path '%s': min/max validation requires a number, got %T", group.Path, val)
			}
			if rules.Min != nil && num < *rules.Min {
				return fmt.Errorf("path '%s': value %v is less than minimum %v", group.Path, num, *rules.Min)
			}
			if rules.Max != nil && num > *rules.Max {
				return fmt.Errorf("path '%s': value %v is greater than maximum %v", group.Path, num, *rules.Max)
			}
		}

		// MinLength for strings
		if rules.MinLength != nil {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("path '%s': minLength validation requires a string, got %T", group.Path, val)
			}
			if len(str) < *rules.MinLength {
				return fmt.Errorf("path '%s': length %d is less than minimum length %d", group.Path, len(str), *rules.MinLength)
			}
		}

		// Enum for strings
		if len(rules.Enum) > 0 {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("path '%s': enum validation requires a string, got %T", group.Path, val)
			}
			match := false
			for _, enumVal := range rules.Enum {
				if str == enumVal {
					match = true
					break
				}
			}
			if !match {
				return fmt.Errorf("path '%s': value '%s' is not in the allowed list %v", group.Path, str, rules.Enum)
			}
		}

		// Regex for strings
		if rules.Regex != "" {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("path '%s': regex validation requires a string, got %T", group.Path, val)
			}
			compiledRegex, err := regexp.Compile(rules.Regex)
			if err != nil {
				return fmt.Errorf("path '%s': invalid regex pattern '%s': %w", group.Path, rules.Regex, err)
			}
			if !compiledRegex.MatchString(str) {
				return fmt.Errorf("path '%s': value '%s' does not match regex pattern '%s'", group.Path, str, rules.Regex)
			}
		}
	}
	return nil
}
