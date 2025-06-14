package schema

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/parser"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// --- Struct Definitions for the Schema File ---

// VarDef defines a variable that can be used for substitution.
type VarDef struct {
	Name         string `yaml:"name" json:"name"`
	Value        string `yaml:"value,omitempty" json:"value,omitempty"`
	FromEnv      string `yaml:"fromEnv,omitempty" json:"fromEnv,omitempty"`
	FromPath     string `yaml:"fromPath,omitempty" json:"fromPath,omitempty"`
	DefaultValue string `yaml:"defaultValue,omitempty" json:"defaultValue,omitempty"`
}

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
	Vars         []VarDef          `yaml:"vars"`
	Generators   []GeneratorDef    `yaml:"generators"`
	Transforms   []TransformDef    `yaml:"transform"`
	Validate     []ValidationGroup `yaml:"validate"`
}

type Ref struct {
	Path   string `yaml:"path"`
	Strict bool   `yaml:"strict"`
}

type GeneratorDef struct {
	Type       string            `yaml:"type"`
	TargetPath string            `yaml:"targetPath"`
	Format     string            `yaml:"format"`
	Sources    map[string]string `yaml:"sources"`
}

type TransformDef struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Case   string `yaml:"case"`
	Prefix string `yaml:"prefix"`
	Value  any    `yaml:"value"`
}

type ValidationGroup struct {
	Path  string         `yaml:"path"`
	Rules ValidationRule `yaml:"rules"`
}

type ValidationRule struct {
	Required  bool     `yaml:"required"`
	Type      string   `yaml:"type"`
	Min       *float64 `yaml:"min"`
	Max       *float64 `yaml:"max"`
	MinLength *int     `yaml:"minLength"`
	Enum      []string `yaml:"enum"`
	Regex     string   `yaml:"regex"`
}

func Load(path string) (*Schema, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "json", "yaml", "yml", "toml":
		break
	case "ini", "env":
		return nil, fmt.Errorf(
			"unsuitable schema format '%s': schema files must be in JSON, YAML, or TOML to support complex structures",
			ext,
		)
	default:
		return nil, fmt.Errorf("unsupported schema file extension: %s", ext)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", path, err)
	}

	data, err := parser.Parse(path, content, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema file %s: %w", path, err)
	}

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to internally process schema data: %w", err)
	}

	var schema Schema
	if err := yaml.Unmarshal(yamlBytes, &schema); err != nil {
		return nil, fmt.Errorf("failed to decode schema structure from %s: %w", path, err)
	}
	return &schema, nil
}

// Process orchestrates the entire schema-driven pipeline with the new steps.
func Process(config map[string]interface{}, schema *Schema, varsFromFile map[string]interface{}, envVars map[string]string) (map[string]interface{}, error) {
	logger.Log("Applying schema...")

	if schema.InputSchema != nil {
		logger.Log("Validating against input schema: %s", schema.InputSchema.Path)
		if err := validateInputSchema(config, schema.InputSchema); err != nil {
			return nil, fmt.Errorf("input schema validation failed: %w", err)
		}
	}

	// 1. Resolve variables, now including envVars
	resolver, err := NewResolver(envVars, varsFromFile, schema.Vars, config)
	if err != nil {
		return nil, fmt.Errorf("variable resolution failed: %w", err)
	}

	// 2. Run generators
	if err := ApplyGenerators(config, schema.Generators, resolver); err != nil {
		return nil, fmt.Errorf("generator failed: %w", err)
	}

	// 3. Run transformers
	if err := ApplyTransforms(config, schema.Transforms, resolver); err != nil {
		return nil, fmt.Errorf("transform failed: %w", err)
	}

	// 4. Substitute variables throughout the config
	processedConfig := Substitute(config, resolver)

	// 5. Validate the final configuration
	if err := ApplyValidations(processedConfig, schema.Validate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if schema.OutputSchema != nil {
		logger.Log("Filtering output against output schema: %s", schema.OutputSchema.Path)
		processedConfig, err = filterOutputSchema(processedConfig, schema.OutputSchema)
		if err != nil {
			return nil, fmt.Errorf("output schema filtering failed: %w", err)
		}
	}

	logger.Log("Schema applied successfully.")
	return processedConfig, nil
}

func validateInputSchema(config map[string]interface{}, ref *Ref) error {
	schemaContent, err := os.ReadFile(ref.Path)
	if err != nil {
		return err
	}
	schemaMap, err := parser.Parse(ref.Path, schemaContent, "")
	if err != nil {
		return err
	}
	return compareMaps(config, schemaMap, "", ref.Strict)
}

func compareMaps(data, schema map[string]interface{}, path string, strict bool) error {
	for key, schemaVal := range schema {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}
		dataVal, ok := data[key]
		if !ok {
			return fmt.Errorf("path '%s': missing required key", currentPath)
		}
		schemaType := reflect.TypeOf(schemaVal)
		dataType := reflect.TypeOf(dataVal)
		if schemaMap, sOK := schemaVal.(map[string]interface{}); sOK {
			if dataMap, dOK := dataVal.(map[string]interface{}); dOK {
				if err := compareMaps(dataMap, schemaMap, currentPath, strict); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("path '%s': type mismatch, expected map, got %v", currentPath, dataType)
			}
		} else if schemaType != dataType {
			if !(schemaType.Kind() == reflect.Int && dataType.Kind() == reflect.Float64) {
				return fmt.Errorf("path '%s': type mismatch, expected %v, got %v", currentPath, schemaType, dataType)
			}
		}
	}
	if strict {
		for key := range data {
			if _, ok := schema[key]; !ok {
				currentPath := key
				if path != "" {
					currentPath = path + "." + key
				}
				return fmt.Errorf("path '%s': unexpected key found in strict mode", currentPath)
			}
		}
	}
	return nil
}

func filterOutputSchema(config map[string]interface{}, ref *Ref) (map[string]interface{}, error) {
	schemaContent, err := os.ReadFile(ref.Path)
	if err != nil {
		return nil, err
	}
	schemaMap, err := parser.Parse(ref.Path, schemaContent, "")
	if err != nil {
		return nil, err
	}
	return projectMap(config, schemaMap, "", ref.Strict)
}

func projectMap(data, schema map[string]interface{}, path string, strict bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, schemaVal := range schema {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}
		if dataVal, ok := data[key]; ok {
			if schemaMap, sOK := schemaVal.(map[string]interface{}); sOK {
				if dataMap, dOK := dataVal.(map[string]interface{}); dOK {
					nestedResult, err := projectMap(dataMap, schemaMap, currentPath, strict)
					if err != nil {
						return nil, err
					}
					result[key] = nestedResult
				} else if strict {
					// Data has a non-map type where schema expects a map, and it's strict mode for outputSchema.
					// This case might be debatable: if schema defines a map, should data also be a map?
					// For now, let's consider it a mismatch if strict is on and types differ like this.
					return nil, fmt.Errorf("outputSchema strict: path '%s' in data is type %T, but schema expects a map", currentPath, dataVal)
				}
				// If not strict, and data is not a map but schema expects one, we just don't include this key.
			} else {
				result[key] = dataVal
			}
		} else if strict {
			// Key is in outputSchema but not in data, and it's strict mode.
			// This implies the field was expected in the output.
			return nil, fmt.Errorf("outputSchema strict: path '%s' defined in output schema but not found in processed configuration", currentPath)
		}
	}

	if strict {
		// Check for keys in data that are NOT in the schema
		for key := range data {
			if _, schemaHasKey := schema[key]; !schemaHasKey {
				currentDataPath := key
				if path != "" {
					currentDataPath = path + "." + key
				}
				return nil, fmt.Errorf("outputSchema strict: path '%s' found in processed configuration but not defined in output schema", currentDataPath)
			}
		}
	}
	return result, nil
}
