package util

import (
	"strconv"
	"strings"
)

// InferType attempts to convert a string value to its most appropriate Go type.
// This is used for environment variable type conversion.
func InferType(value string) interface{} {
	// Handle empty string
	if value == "" {
		return ""
	}
	
	// Check for boolean values (case-insensitive)
	lowerValue := strings.ToLower(value)
	if lowerValue == "true" {
		return true
	}
	if lowerValue == "false" {
		return false
	}
	
	// Try to parse as integer
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		// If it fits in a regular int, return that, otherwise int64
		if intVal >= -2147483648 && intVal <= 2147483647 {
			return int(intVal)
		}
		return intVal
	}
	
	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	
	// Return as string if no other type matches
	return value
}

// TryConvertType attempts to convert a string value to a specific target type.
// This is used when we have schema type information available.
func TryConvertType(value string, targetType string) (interface{}, bool) {
	switch targetType {
	case "bool", "boolean":
		lowerValue := strings.ToLower(value)
		if lowerValue == "true" {
			return true, true
		}
		if lowerValue == "false" {
			return false, true
		}
		return nil, false
		
	case "number", "int", "integer":
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			if intVal >= -2147483648 && intVal <= 2147483647 {
				return int(intVal), true
			}
			return intVal, true
		}
		return nil, false
		
	case "float", "float64", "double":
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal, true
		}
		return nil, false
		
	case "string":
		return value, true
		
	default:
		// Unknown target type, try automatic inference
		return InferType(value), true
	}
}
