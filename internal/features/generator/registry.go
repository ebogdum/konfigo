package generator

import (
	"fmt"
	"konfigo/internal/logger"
)

// Apply applies all generators to the configuration using the default registry.
func Apply(config map[string]interface{}, definitions []Definition, resolver VariableResolver) error {
	if len(definitions) == 0 {
		return nil
	}

	logger.Debug("Applying %d generators...", len(definitions))

	registry := NewRegistry()

	if err := ValidateDefinitionsWithRegistry(definitions, registry); err != nil {
		return err
	}

	for _, def := range definitions {
		logger.Debug("  - Processing generator type '%s'", def.Type)

		generator, exists := registry.Get(def.Type)
		if !exists {
			return fmt.Errorf("unsupported generator type: %s", def.Type)
		}

		if err := generator.Generate(config, def, resolver); err != nil {
			return fmt.Errorf("generator '%s' failed: %w", def.Type, err)
		}
	}


	logger.Debug("All generators applied successfully")
	return nil
}

// ApplyWithRegistry applies generators using a custom registry.
func ApplyWithRegistry(config map[string]interface{}, definitions []Definition, resolver VariableResolver, registry Registry) error {
	if len(definitions) == 0 {
		return nil
	}

	logger.Debug("Applying %d generators with custom registry...", len(definitions))

	if err := ValidateDefinitionsWithRegistry(definitions, registry); err != nil {
		return err
	}

	for _, def := range definitions {
		generator, exists := registry.Get(def.Type)
		if !exists {
			return fmt.Errorf("unsupported generator type: %s", def.Type)
		}

		if err := generator.Generate(config, def, resolver); err != nil {
			return fmt.Errorf("generator '%s' failed: %w", def.Type, err)
		}
	}


	return nil
}

// ValidateDefinitions validates all generator definitions using a new registry.
func ValidateDefinitions(definitions []Definition) error {
	return ValidateDefinitionsWithRegistry(definitions, NewRegistry())
}

// ValidateDefinitionsWithRegistry validates all generator definitions using the provided registry.
func ValidateDefinitionsWithRegistry(definitions []Definition, registry Registry) error {
	for i, def := range definitions {
		generator, exists := registry.Get(def.Type)
		if !exists {
			return fmt.Errorf("definition %d: unsupported generator type: %s", i, def.Type)
		}

		if validator, ok := generator.(interface {
			ValidateDefinition(Definition) error
		}); ok {
			if err := validator.ValidateDefinition(def); err != nil {
				return fmt.Errorf("definition %d: %w", i, err)
			}
		}
	}

	return nil
}
