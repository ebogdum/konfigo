package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"
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

	// Check if destination path exists — warn if not, as this creates a new key
	if _, exists := util.GetNestedValue(config, def.Path); !exists {
		logger.Warn("replaceKey: destination path '%s' does not exist, creating new key", def.Path)
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

	// Same overlap hazard as renameKey: this writes to 'path' and then deletes
	// 'target', so when one is nested inside the other the delete removes the
	// value that was just written, or the write clobbers the parent holding it.
	// Either way data disappears with no error, so refuse the definition.
	if strings.HasPrefix(def.Path, def.Target+".") {
		return fmt.Errorf("replaceKey transformer: 'path' %q is nested inside 'target' %q; deleting the target would discard the replaced value", def.Path, def.Target)
	}
	if strings.HasPrefix(def.Target, def.Path+".") {
		return fmt.Errorf("replaceKey transformer: 'target' %q is nested inside 'path' %q; writing to the parent would drop the target's siblings", def.Target, def.Path)
	}

	return nil
}
