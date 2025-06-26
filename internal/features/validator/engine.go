package validator

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// Apply applies all validation groups to the configuration.
func Apply(config map[string]interface{}, groups []Group) error {
	if len(groups) == 0 {
		return nil
	}
	
	logger.Debug("Applying %d validation groups...", len(groups))
	
	registry := NewRegistry()
	
	for _, group := range groups {
		val, found := util.GetNestedValue(config, group.Path)
		
		// Check required first
		if group.Rules.Required && !found {
			return fmt.Errorf("path '%s' is required but not found", group.Path)
		}
		
		// Skip other validations if not found and not required
		if !found {
			continue
		}
		
		logger.Debug("  - Validating path '%s'", group.Path)
		
		// Apply each validation rule
		for _, validator := range registry.GetValidators() {
			if err := validator.Validate(val, group.Path, group.Rules); err != nil {
				return err
			}
		}
	}
	
	return nil
}
