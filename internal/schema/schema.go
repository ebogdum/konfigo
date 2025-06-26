// Package schema provides schema loading, validation, and processing capabilities.
//
// This package handles loading schema files and processing configuration data
// through various features defined in the schema including:
// - Configuration validation against schema rules
// - Variable resolution and substitution
// - Configuration generation and transformation
// - Batch processing directives (konfigo_forEach)
//
// Schema Structure:
//
//	A schema defines the expected structure of configuration data and
//	includes processing directives for features like validation, generation,
//	transformation, and variable resolution.
//
// Usage:
//
//	schema, err := schema.LoadFromFile("config.schema.yaml")
//	if err != nil {
//	    return err
//	}
//	processed, err := schema.Process(configData)
package schema

import (
	"konfigo/internal/errors"
	"konfigo/internal/features/generator"
	"konfigo/internal/features/transformer"
	"konfigo/internal/features/validator"
	"konfigo/internal/features/variables"
	"konfigo/internal/parser"
	"konfigo/internal/reader"

	"gopkg.in/yaml.v3"
)

// --- Struct Definitions for the Schema File ---

// KonfigoForEachOutput defines the output configuration for batch processing.
type KonfigoForEachOutput struct {
	FilenamePattern string `yaml:"filenamePattern" json:"filenamePattern"`
	Format          string `yaml:"format,omitempty" json:"format,omitempty"` // e.g., json, yaml, toml
}

// KonfigoForEach defines the structure for batch processing directives.
// It will be looked for in the primary variables file.
type KonfigoForEach struct {
	Items     []map[string]interface{} `yaml:"items,omitempty" json:"items,omitempty"`
	ItemFiles []string                 `yaml:"itemFiles,omitempty" json:"itemFiles,omitempty"`
	Output    KonfigoForEachOutput     `yaml:"output" json:"output"`
	// GlobalVars will hold variables defined outside konfigo_forEach in the main vars file
	GlobalVars map[string]interface{} `yaml:"-" json:"-"` // Loaded separately
}

// Schema represents the entire Konfigo schema structure.
type Schema struct {
	APIVersion   string            `yaml:"apiVersion"`
	InputSchema  *Ref              `yaml:"inputSchema"`
	OutputSchema *Ref              `yaml:"outputSchema"`
	Immutable    []string          `yaml:"immutable"`
	Vars         []variables.Definition `yaml:"vars"`
	Generators   []generator.Definition `yaml:"generators"`
	Transforms   []transformer.Definition `yaml:"transform"`
	Validate     []validator.Group `yaml:"validate"`
}

// Ref represents a reference to another schema file.
type Ref struct {
	Path   string `yaml:"path"`
	Strict bool   `yaml:"strict"`
}

// Load loads and parses a schema file from the given path.
func Load(path string) (*Schema, error) {
	format := parser.DetectFormat(path)
	
	// Validate that the format is suitable for schema files
	if err := parser.ValidateSchemaFormat(format); err != nil {
		return nil, err
	}

	content, err := reader.ReadFile(path)
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeSchemaLoad, "failed to read schema file", err).WithContext("file", path)
	}

	data, err := parser.Parse(path, content, "")
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeParsing, "failed to parse schema file", err).WithContext("file", path)
	}

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeSchemaProcess, "failed to internally process schema data", err)
	}

	var schema Schema
	if err := yaml.Unmarshal(yamlBytes, &schema); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeSchemaLoad, "failed to decode schema structure", err).WithContext("file", path)
	}
	return &schema, nil
}

// Process is a convenience function that creates a processor and processes the configuration.
// This maintains backward compatibility while allowing for more flexible processing.
func Process(config map[string]interface{}, schema *Schema, varsFromFile map[string]interface{}, envVars map[string]string) (map[string]interface{}, error) {
	processor := NewProcessor()
	return processor.Process(config, schema, varsFromFile, envVars)
}
