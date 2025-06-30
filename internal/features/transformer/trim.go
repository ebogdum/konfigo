package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"
)

// TrimType is the type identifier for the trim transformer.
const TrimType = "trim"

// TrimTransformer trims whitespace or specified pattern from string values.
type TrimTransformer struct{}

// Type returns the transformer type.
func (t *TrimTransformer) Type() string {
	return TrimType
}

// Transform implements the trim transformation logic.
// It trims the specified pattern or whitespace from a string value.
func (t *TrimTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying trim transform at path '%s'", def.Path)
	
	// Get the value from the specified path
	value, found := util.GetNestedValue(config, def.Path)
	if !found {
		return fmt.Errorf("trim: path '%s' not found", def.Path)
	}
	
	// Ensure the value is a string
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("trim: value at path '%s' is not a string (got %T)", def.Path, value)
	}
	
	// Apply trimming
	var newValue string
	if def.Pattern == "" {
		// Default: trim whitespace
		newValue = strings.TrimSpace(strValue)
	} else {
		// Trim specified pattern
		newValue = strings.Trim(strValue, def.Pattern)
	}
	
	// Set the trimmed value
	util.SetNestedValue(config, def.Path, newValue)
	
	logger.Debug("    Trimmed value from '%s' to '%s'", strValue, newValue)
	return nil
}

// ValidateDefinition validates a trim transformer definition.
func (t *TrimTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("trim transformer: 'path' is required")
	}
	
	return nil
}
