package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// AddKeyPrefixType is the type identifier for the add key prefix transformer.
const AddKeyPrefixType = "addKeyPrefix"

// AddKeyPrefixTransformer adds a prefix to all keys in a map at the specified path.
type AddKeyPrefixTransformer struct{}

// Type returns the transformer type.
func (t *AddKeyPrefixTransformer) Type() string {
	return AddKeyPrefixType
}

// Transform implements the add key prefix transformation logic.
// It adds a prefix to all keys in a map structure.
func (t *AddKeyPrefixTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying addKeyPrefix transform at path '%s' with prefix '%s'", def.Path, def.Prefix)
	
	// Get the value from the specified path
	value, found := util.GetNestedValue(config, def.Path)
	if !found {
		return fmt.Errorf("addKeyPrefix: path '%s' not found", def.Path)
	}
	
	// Ensure the value is a map
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("addKeyPrefix: value at path '%s' is not a map (got %T)", def.Path, value)
	}
	
	// Create a new map with prefixed keys
	newMap := make(map[string]interface{})
	for key, val := range mapValue {
		newKey := def.Prefix + key
		newMap[newKey] = val
	}
	
	// Set the new map at the path
	util.SetNestedValue(config, def.Path, newMap)
	
	logger.Debug("    Added prefix '%s' to %d keys", def.Prefix, len(mapValue))
	return nil
}

// ValidateDefinition validates an add key prefix transformer definition.
func (t *AddKeyPrefixTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("addKeyPrefix transformer: 'path' is required")
	}
	
	if def.Prefix == "" {
		return fmt.Errorf("addKeyPrefix transformer: 'prefix' is required")
	}
	
	return nil
}
