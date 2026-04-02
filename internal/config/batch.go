package config

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/schema"

	"gopkg.in/yaml.v3"
)

// ExtractForEachFromVars extracts forEach directive from variables file.
func ExtractForEachFromVars(varsFromFile map[string]interface{}) (*schema.KonfigoForEach, map[string]interface{}, error) {
	if varsFromFile == nil {
		return nil, nil, nil
	}

	var forEachConfig *schema.KonfigoForEach
	globalVars := make(map[string]interface{})

	// Separate forEach from other global vars
	for k, v := range varsFromFile {
		if k == "forEach" {
			// Convert to KonfigoForEach struct
			yamlBytes, err := yaml.Marshal(v)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal forEach directive: %w", err)
			}

			const maxForEachSize = 10 * 1024 * 1024 // 10 MiB
			if len(yamlBytes) > maxForEachSize {
				return nil, nil, fmt.Errorf("forEach directive exceeds maximum allowed size of %d bytes", maxForEachSize)
			}
			forEachConfig = &schema.KonfigoForEach{}
			if err := yaml.Unmarshal(yamlBytes, forEachConfig); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal forEach directive: %w", err)
			}

			logger.Debug("Found forEach directive.")
		} else {
			globalVars[k] = v
		}
	}

	if forEachConfig != nil {
		forEachConfig.GlobalVars = globalVars
	}

	return forEachConfig, globalVars, nil
}
