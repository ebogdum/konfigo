package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DetectFormat detects the file format from the file path extension.
func DetectFormat(filePath string) string {
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	return strings.ToLower(ext)
}

// IsFormatSupported checks if the given format is supported.
func IsFormatSupported(format string) bool {
	switch strings.ToLower(format) {
	case "json", "yaml", "yml", "toml", "ini", "env":
		return true
	default:
		return false
	}
}

// IsSchemaFormat checks if the given format is suitable for schema files.
// Only JSON, YAML, and TOML are suitable for schema files due to their
// support for complex nested structures.
func IsSchemaFormat(format string) bool {
	switch strings.ToLower(format) {
	case "json", "yaml", "yml", "toml":
		return true
	default:
		return false
	}
}

// ValidateSchemaFormat validates that a format is suitable for schema files.
func ValidateSchemaFormat(format string) error {
	if !IsSchemaFormat(format) {
		return fmt.Errorf(
			"unsuitable schema format '%s': schema files must be in JSON, YAML, or TOML to support complex structures",
			format,
		)
	}
	return nil
}

// NormalizeFormat normalizes format names (e.g., "yml" -> "yaml").
func NormalizeFormat(format string) string {
	switch strings.ToLower(format) {
	case "yml":
		return "yaml"
	default:
		return strings.ToLower(format)
	}
}
