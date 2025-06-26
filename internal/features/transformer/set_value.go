package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// SetValueType is the type identifier for the set value transformer.
const SetValueType = "setValue"

// SetValueTransformer sets a specific value at a configuration path.
type SetValueTransformer struct{}

// Type returns the transformer type.
func (t *SetValueTransformer) Type() string {
	return SetValueType
}

// Transform implements the set value transformation logic.
// It sets a specific value at the specified path, creating the path if necessary.
func (t *SetValueTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying setValue transform at path '%s'", def.Path)
	
	// Set the value at the specified path
	util.SetNestedValue(config, def.Path, def.Value)
	
	logger.Debug("    Set value '%v' at path '%s'", def.Value, def.Path)
	return nil
}

// ValidateDefinition validates a set value transformer definition.
func (t *SetValueTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("setValue transformer: 'path' is required")
	}
	
	// Note: Value can be nil, so we don't validate it
	
	return nil
}
