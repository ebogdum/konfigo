package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// RenameKeyType is the type identifier for the rename key transformer.
const RenameKeyType = "renameKey"

// RenameKeyTransformer renames configuration keys by moving values from one path to another.
type RenameKeyTransformer struct{}

// Type returns the transformer type.
func (t *RenameKeyTransformer) Type() string {
	return RenameKeyType
}

// Transform implements the rename key transformation logic.
// It moves a value from the 'From' path to the 'To' path and removes the original.
func (t *RenameKeyTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying renameKey transform from '%s' to '%s'", def.From, def.To)
	
	// Get the value from the source path
	value, found := util.GetNestedValue(config, def.From)
	if !found {
		return fmt.Errorf("renameKey: source path '%s' not found", def.From)
	}
	
	// Set the value at the destination path
	util.SetNestedValue(config, def.To, value)
	
	// Remove the value from the source path
	util.DeleteNestedValue(config, def.From)
	
	logger.Debug("    Renamed key from '%s' to '%s'", def.From, def.To)
	return nil
}

// ValidateDefinition validates a rename key transformer definition.
func (t *RenameKeyTransformer) ValidateDefinition(def Definition) error {
	if def.From == "" {
		return fmt.Errorf("renameKey transformer: 'from' path is required")
	}
	
	if def.To == "" {
		return fmt.Errorf("renameKey transformer: 'to' path is required")
	}
	
	if def.From == def.To {
		return fmt.Errorf("renameKey transformer: 'from' and 'to' paths cannot be the same")
	}
	
	return nil
}
