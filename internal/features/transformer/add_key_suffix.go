package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// AddKeySuffixType is the type identifier for the add key suffix transformer.
const AddKeySuffixType = "addKeySuffix"

// AddKeySuffixTransformer adds a suffix to all keys in a map at the specified path.
type AddKeySuffixTransformer struct{}

// Type returns the transformer type.
func (t *AddKeySuffixTransformer) Type() string {
	return AddKeySuffixType
}

// Transform implements the add key suffix transformation logic.
// It adds a suffix to all keys in a map structure.
func (t *AddKeySuffixTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying addKeySuffix transform at path '%s' with suffix '%s'", def.Path, def.Suffix)
	
	// Get the value from the specified path
	value, found := util.GetNestedValue(config, def.Path)
	if !found {
		return fmt.Errorf("addKeySuffix: path '%s' not found", def.Path)
	}
	
	// Ensure the value is a map
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("addKeySuffix: value at path '%s' is not a map (got %T)", def.Path, value)
	}
	
	// Create a new map with suffixed keys
	newMap := make(map[string]interface{})
	for key, val := range mapValue {
		newKey := key + def.Suffix
		newMap[newKey] = val
	}
	
	// Set the new map at the path
	util.SetNestedValue(config, def.Path, newMap)
	
	logger.Debug("    Added suffix '%s' to %d keys", def.Suffix, len(mapValue))
	return nil
}

// ValidateDefinition validates an add key suffix transformer definition.
func (t *AddKeySuffixTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("addKeySuffix transformer: 'path' is required")
	}
	
	if def.Suffix == "" {
		return fmt.Errorf("addKeySuffix transformer: 'suffix' is required")
	}
	
	return nil
}
