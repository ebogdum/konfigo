package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// ReplaceKeyType is the type identifier for the replace key transformer.
const ReplaceKeyType = "replaceKey"

// ReplaceKeyTransformer replaces a key with value from target, then deletes the target.
type ReplaceKeyTransformer struct{}

// Type returns the transformer type.
func (t *ReplaceKeyTransformer) Type() string {
	return ReplaceKeyType
}

// Transform implements the replace key transformation logic.
// It takes the value from the target path and replaces the value at the main path, then deletes the target.
func (t *ReplaceKeyTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying replaceKey transform at path '%s' with target '%s'", def.Path, def.Target)
	
	// Get the value from the target path
	targetValue, found := util.GetNestedValue(config, def.Target)
	if !found {
		return fmt.Errorf("replaceKey: target path '%s' not found", def.Target)
	}
	
	// Set the target value at the main path
	util.SetNestedValue(config, def.Path, targetValue)
	
	// Delete the target path
	util.DeleteNestedValue(config, def.Target)
	
	logger.Debug("    Replaced value at '%s' with value from '%s' and deleted target", def.Path, def.Target)
	return nil
}

// ValidateDefinition validates a replace key transformer definition.
func (t *ReplaceKeyTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("replaceKey transformer: 'path' is required")
	}
	
	if def.Target == "" {
		return fmt.Errorf("replaceKey transformer: 'target' is required")
	}
	
	if def.Path == def.Target {
		return fmt.Errorf("replaceKey transformer: 'path' and 'target' cannot be the same")
	}
	
	return nil
}
