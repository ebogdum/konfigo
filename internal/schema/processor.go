package schema

import (
	"fmt"
	"konfigo/internal/errors"
	"konfigo/internal/features/generator"
	"konfigo/internal/features/input_schema"
	"konfigo/internal/features/transformer"
	"konfigo/internal/features/validator"
	"konfigo/internal/features/variables"
	"konfigo/internal/logger"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
	"konfigo/internal/util"
	"path/filepath"
	"reflect"
)

// Processor handles the orchestration of schema-driven configuration processing.
type Processor struct{}

// NewProcessor creates a new processor instance.
func NewProcessor() *Processor {
	return &Processor{}
}

// Process orchestrates the entire schema-driven pipeline with the new steps.
func (p *Processor) Process(config map[string]interface{}, schema *Schema, varsFromFile map[string]interface{}, envVars map[string]string) (map[string]interface{}, error) {
	logger.Log("Applying schema...")

	if schema.InputSchema != nil {
		resolvedRef, err := resolveRefPath(schema.BaseDir, schema.InputSchema)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeSchemaLoad, "invalid input schema ref", err)
		}
		if err := input_schema.Validate(config, (*input_schema.Ref)(resolvedRef)); err != nil {
			return nil, errors.WrapError(errors.ErrorTypeValidation, "input schema validation failed", err)
		}
	}

	// Build immutable paths set for enforcement during generation/transformation
	immutableSet := make(map[string]struct{}, len(schema.Immutable))
	for _, ip := range schema.Immutable {
		immutableSet[ip] = struct{}{}
	}

	// 1. Resolve variables, now including envVars
	resolver, err := variables.NewResolver(envVars, varsFromFile, schema.Vars, config)
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeVarResolution, "variable resolution failed", err)
	}

	// Snapshot immutable values before generators/transformers
	immutableSnapshot := make(map[string]interface{}, len(immutableSet))
	for ip := range immutableSet {
		if val, found := util.GetNestedValue(config, ip); found {
			immutableSnapshot[ip] = val
		}
	}

	// 2. Run generators
	if err := generator.Apply(config, schema.Generators, resolver); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeInternal, "generator failed", err)
	}

	// 3. Run transformers
	if err := p.applyTransforms(config, schema.Transforms, resolver); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeInternal, "transform failed", err)
	}

	// Restore immutable values if they were modified by generators/transformers
	for ip, originalVal := range immutableSnapshot {
		currentVal, found := util.GetNestedValue(config, ip)
		if !found || !reflect.DeepEqual(currentVal, originalVal) {
			logger.Warn("Restoring immutable path '%s' that was modified by generator/transformer", ip)
			util.SetNestedValue(config, ip, originalVal)
		}
	}

	// 4. Substitute variables throughout the config
	processedConfig := variables.Substitute(config, resolver)

	// 5. Validate the final configuration
	if err := validator.Apply(processedConfig, schema.Validate); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeValidation, "validation failed", err)
	}

	if schema.OutputSchema != nil {
		resolvedOutRef, err := resolveRefPath(schema.BaseDir, schema.OutputSchema)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeSchemaLoad, "invalid output schema ref", err)
		}
		logger.Log("Filtering output against output schema: %s", resolvedOutRef.Path)
		processedConfig, err = p.filterOutputSchema(processedConfig, resolvedOutRef)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeValidation, "output schema filtering failed", err)
		}
	}

	logger.Log("Schema applied successfully.")
	return processedConfig, nil
}

// applyTransforms applies a list of transformer definitions to the configuration.
func (p *Processor) applyTransforms(config map[string]interface{}, transforms []transformer.Definition, resolver variables.Resolver) error {
	return transformer.Apply(config, transforms, resolver)
}

// filterOutputSchema filters the configuration against an output schema.
func (p *Processor) filterOutputSchema(config map[string]interface{}, ref *Ref) (map[string]interface{}, error) {
	schemaContent, err := reader.ReadFile(ref.Path)
	if err != nil {
		return nil, err
	}
	schemaMap, err := parser.Parse(ref.Path, schemaContent, "")
	if err != nil {
		return nil, err
	}
	return p.projectMap(config, schemaMap, "", ref.Strict)
}

// projectMap projects data onto a schema map, optionally enforcing strict validation.
func (p *Processor) projectMap(data, schema map[string]interface{}, path string, strict bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, schemaVal := range schema {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}
		if dataVal, ok := data[key]; ok {
			if schemaMap, sOK := schemaVal.(map[string]interface{}); sOK {
				if dataMap, dOK := dataVal.(map[string]interface{}); dOK {
					nestedResult, err := p.projectMap(dataMap, schemaMap, currentPath, strict)
					if err != nil {
						return nil, err
					}
					result[key] = nestedResult
				} else if strict {
					// Data has a non-map type where schema expects a map, and it's strict mode for outputSchema.
					// This case might be debatable: if schema defines a map, should data also be a map?
					// For now, let's consider it a mismatch if strict is on and types differ like this.
					return nil, errors.NewErrorf(errors.ErrorTypeValidation, "outputSchema strict: path '%s' in data is type %T, but schema expects a map", currentPath, dataVal)
				}
				// If not strict, and data is not a map but schema expects one, we just don't include this key.
			} else {
				result[key] = dataVal
			}
		} else if strict {
			// Key is in outputSchema but not in data, and it's strict mode.
			// This implies the field was expected in the output.
			return nil, errors.NewErrorf(errors.ErrorTypeValidation, "outputSchema strict: path '%s' defined in output schema but not found in processed configuration", currentPath)
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
				return nil, errors.NewErrorf(errors.ErrorTypeValidation, "outputSchema strict: path '%s' found in processed configuration but not defined in output schema", currentDataPath)
			}
		}
	}
	return result, nil
}

// resolveRefPath resolves a Ref's path relative to the schema's base directory.
// Since schemas are provided by the user (not untrusted input), path traversal
// is allowed. The function cleans the path to prevent accidental issues.
func resolveRefPath(baseDir string, ref *Ref) (*Ref, error) {
	if ref == nil || ref.Path == "" {
		return ref, nil
	}

	resolved := ref.Path
	if !filepath.IsAbs(resolved) && baseDir != "" {
		resolved = filepath.Join(baseDir, resolved)
	}

	absResolved, err := filepath.Abs(resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve schema ref path %q: %w", ref.Path, err)
	}

	return &Ref{
		Path:   absResolved,
		Strict: ref.Strict,
	}, nil
}
