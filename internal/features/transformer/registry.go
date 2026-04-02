package transformer

import (
	"fmt"
	"konfigo/internal/logger"
)

// Apply applies all transformations to the configuration using the default registry.
func Apply(config map[string]interface{}, definitions []Definition, resolver VariableResolver) error {
	if len(definitions) == 0 {
		return nil
	}

	logger.Debug("Applying %d transformations...", len(definitions))

	registry := NewRegistry()

	for _, def := range definitions {
		// Substitute variables first, then validate the resolved definition
		processedDef := substituteInDefinition(def, resolver)

		if err := validateSingleDefinition(registry, processedDef); err != nil {
			return err
		}

		logger.Debug("  - Processing transformer type '%s'", processedDef.Type)

		transformer, exists := registry.Get(processedDef.Type)
		if !exists {
			return fmt.Errorf("unsupported transformer type: %s", processedDef.Type)
		}

		if err := transformer.Transform(config, processedDef); err != nil {
			return fmt.Errorf("transformer '%s' failed: %w", processedDef.Type, err)
		}
	}

	logger.Debug("All transformations applied successfully")
	return nil
}

// ApplyWithRegistry applies transformations using a custom registry.
func ApplyWithRegistry(config map[string]interface{}, definitions []Definition, resolver VariableResolver, registry Registry) error {
	if len(definitions) == 0 {
		return nil
	}

	logger.Debug("Applying %d transformations with custom registry...", len(definitions))

	for _, def := range definitions {
		processedDef := substituteInDefinition(def, resolver)

		if err := validateSingleDefinition(registry, processedDef); err != nil {
			return err
		}

		transformer, exists := registry.Get(processedDef.Type)
		if !exists {
			return fmt.Errorf("unsupported transformer type: %s", processedDef.Type)
		}

		if err := transformer.Transform(config, processedDef); err != nil {
			return fmt.Errorf("transformer '%s' failed: %w", processedDef.Type, err)
		}
	}

	return nil
}

// validateSingleDefinition validates a single transformer definition against the registry.
func validateSingleDefinition(registry Registry, def Definition) error {
	transformer, exists := registry.Get(def.Type)
	if !exists {
		return fmt.Errorf("unsupported transformer type: %s", def.Type)
	}
	if validator, ok := transformer.(interface {
		ValidateDefinition(Definition) error
	}); ok {
		if err := validator.ValidateDefinition(def); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDefinitions validates all transformer definitions (pre-substitution, for structural checks only).
func ValidateDefinitions(definitions []Definition) error {
	registry := NewRegistry()

	for i, def := range definitions {
		_, exists := registry.Get(def.Type)
		if !exists {
			return fmt.Errorf("definition %d: unsupported transformer type: %s", i, def.Type)
		}
	}

	return nil
}

// substituteInDefinition performs variable substitution on definition fields.
func substituteInDefinition(def Definition, resolver VariableResolver) Definition {
	if resolver == nil {
		return def
	}

	// Create a copy and substitute variables
	processed := def
	processed.Path = resolver.SubstituteString(def.Path)
	processed.From = resolver.SubstituteString(def.From)
	processed.To = resolver.SubstituteString(def.To)
	processed.Prefix = resolver.SubstituteString(def.Prefix)
	processed.Suffix = resolver.SubstituteString(def.Suffix)
	processed.Pattern = resolver.SubstituteString(def.Pattern)
	processed.Target = resolver.SubstituteString(def.Target)
	processed.Case = resolver.SubstituteString(def.Case)

	// Handle Value field if it's a string
	if s, ok := def.Value.(string); ok {
		processed.Value = resolver.SubstituteString(s)
	}

	return processed
}
