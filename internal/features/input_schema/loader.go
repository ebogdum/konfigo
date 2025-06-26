// Package input_schema handles input schema validation for configuration files.
// It provides functionality to validate configuration structure against a schema definition.
package input_schema

import (
	"fmt"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
)

// Ref represents a reference to a schema file with validation settings.
type Ref struct {
	Path   string `yaml:"path" json:"path"`
	Strict bool   `yaml:"strict" json:"strict"`
}

// LoadSchemaMap loads and parses a schema file, returning it as a map structure.
func LoadSchemaMap(ref *Ref) (map[string]interface{}, error) {
	content, err := reader.ReadFile(ref.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input schema file %s: %w", ref.Path, err)
	}
	
	schemaMap, err := parser.Parse(ref.Path, content, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse input schema file %s: %w", ref.Path, err)
	}
	
	return schemaMap, nil
}

// ValidateFormatSupport ensures the schema file is in a supported format for schemas.
func ValidateFormatSupport(filePath string) error {
	if !reader.IsSupported(filePath) {
		return fmt.Errorf("unsupported file format for schema: %s", filePath)
	}
	
	// Additional validation for schema-specific format restrictions can be added here
	// For now, all supported formats are allowed for input schemas
	return nil
}
