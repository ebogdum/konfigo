// Package parser provides configuration file parsing capabilities.
//
// This package handles parsing various configuration file formats into
// standard Go data structures. It uses a registry-based approach to support
// pluggable parsers for different formats.
//
// Supported formats:
// - JSON: Standard JSON format parsing
// - YAML: YAML format with full specification support
// - TOML: TOML format parsing
// - INI: INI/Properties file format
// - ENV: Environment variable format (KEY=value)
//
// Usage:
//
//	content := []byte(`{"key": "value"}`)
//	data, err := parser.Parse(content, "json")
//	if err != nil {
//	    return err
//	}
package parser

import (
	"konfigo/internal/errors"
)

// defaultRegistry is the global registry instance.
var defaultRegistry = NewRegistry()

// Parse takes file content and parses it into a map.
// It uses the formatOverride if provided, otherwise it detects the format
// from the filePath extension.
func Parse(filePath string, content []byte, formatOverride string) (map[string]interface{}, error) {
	format := formatOverride
	if format == "" {
		format = DetectFormat(filePath)
	}
	
	format = NormalizeFormat(format)
	
	parser, exists := defaultRegistry.Get(format)
	if !exists {
		return nil, errors.NewErrorf(errors.ErrorTypeInvalidFormat, "unsupported file format: %s for file %s", format, filePath)
	}
	
	return parser.Parse(content)
}
