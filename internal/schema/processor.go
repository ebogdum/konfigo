package schema

import (
	"konfigo/internal/errors"
	"konfigo/internal/features/generator"
	"konfigo/internal/features/input_schema"
	"konfigo/internal/features/transformer"
	"konfigo/internal/features/validator"
	"konfigo/internal/features/variables"
	"konfigo/internal/logger"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
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
		if err := input_schema.Validate(config, (*input_schema.Ref)(schema.InputSchema)); err != nil {
			return nil, errors.WrapError(errors.ErrorTypeValidation, "input schema validation failed", err)
		}
	}

	// 1. Resolve variables, now including envVars
	resolver, err := variables.NewResolver(envVars, varsFromFile, schema.Vars, config)
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeVarResolution, "variable resolution failed", err)
	}

	// 2. Run generators
	if err := generator.Apply(config, schema.Generators, resolver); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeInternal, "generator failed", err)
	}

	// 3. Run transformers
	if err := p.applyTransforms(config, schema.Transforms, resolver); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeInternal, "transform failed", err)
	}

	// 4. Substitute variables throughout the config
	processedConfig := variables.Substitute(config, resolver)

	// 5. Validate the final configuration
	if err := validator.Apply(processedConfig, schema.Validate); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeValidation, "validation failed", err)
	}

	if schema.OutputSchema != nil {
		logger.Log("Filtering output against output schema: %s", schema.OutputSchema.Path)
		processedConfig, err = p.filterOutputSchema(processedConfig, schema.OutputSchema)
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
