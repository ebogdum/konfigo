package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// DeleteKeyType is the type identifier for the delete key transformer.
const DeleteKeyType = "deleteKey"

// DeleteKeyTransformer deletes a key at the specified path.
type DeleteKeyTransformer struct{}

// Type returns the transformer type.
func (t *DeleteKeyTransformer) Type() string {
	return DeleteKeyType
}

// Transform implements the delete key transformation logic.
// It removes a key from the configuration.
func (t *DeleteKeyTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying deleteKey transform at path '%s'", def.Path)
	
	// Check if the path exists
	_, found := util.GetNestedValue(config, def.Path)
	if !found {
		return fmt.Errorf("deleteKey: path '%s' not found", def.Path)
	}
	
	// Delete the value at the path
	util.DeleteNestedValue(config, def.Path)
	
	logger.Debug("    Deleted key at path '%s'", def.Path)
	return nil
}

// ValidateDefinition validates a delete key transformer definition.
func (t *DeleteKeyTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("deleteKey transformer: 'path' is required")
	}
	
	return nil
}
