// Package marshaller provides configuration data marshalling to various output formats.
//
// This package handles converting processed configuration data into different
// output formats including JSON, YAML, TOML, and environment variables.
// The package uses a registry-based approach to support pluggable marshallers.
//
// Supported formats:
// - JSON: Standard JSON format
// - YAML: YAML format with proper indentation
// - TOML: TOML format with sections
// - ENV: Environment variable format (KEY=value)
//
// Usage:
//
//	data := map[string]interface{}{"key": "value"}
//	output, err := marshaller.Marshal(data, "yaml")
//	if err != nil {
//	    return err
//	}
package marshaller

import (
	"konfigo/internal/errors"
	"strings"
)

// defaultRegistry is the global registry instance.
var defaultRegistry = NewRegistry()

// Marshal takes the final merged data and a format string, returning the
// data as a byte slice in the specified format.
func Marshal(data map[string]interface{}, format string) ([]byte, error) {
	normalizedFormat := strings.ToLower(format)
	if normalizedFormat == "yml" {
		normalizedFormat = "yaml"
	}
	
	marshaller, exists := defaultRegistry.Get(normalizedFormat)
	if !exists {
		return nil, errors.NewErrorf(errors.ErrorTypeInvalidFormat, "unsupported output format: %s", format)
	}
	
	return marshaller.Marshal(data)
}
