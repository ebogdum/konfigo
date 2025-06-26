package config

import (
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"os"
	"strings"
)

// Environment handles loading configuration from environment variables.
type Environment struct {
	keyPrefix string
	varPrefix string
}

// NewEnvironment creates a new environment loader with default prefixes.
func NewEnvironment() *Environment {
	return &Environment{
		keyPrefix: "KONFIGO_KEY_",
		varPrefix: "KONFIGO_VAR_",
	}
}

// NewEnvironmentWithPrefixes creates a new environment loader with custom prefixes.
func NewEnvironmentWithPrefixes(keyPrefix, varPrefix string) *Environment {
	return &Environment{
		keyPrefix: keyPrefix,
		varPrefix: varPrefix,
	}
}

// LoadResult contains the results of loading from environment.
type LoadResult struct {
	Config *Config
	Vars   map[string]string
}

// Load loads configuration and variables from environment variables.
// Returns configuration from KONFIGO_KEY_ variables and variables from KONFIGO_VAR_ variables.
func (e *Environment) Load() *LoadResult {
	config := New()
	vars := make(map[string]string)

	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key, value := parts[0], parts[1]

		if strings.HasPrefix(key, e.keyPrefix) {
			configKey := strings.TrimPrefix(key, e.keyPrefix)
			logger.Debug("  - Loading from env config: %s -> %s", key, configKey)
			
			// Apply type inference to the environment variable value
			typedValue := util.InferType(value)
			logger.Debug("    Type inference: '%s' (%T) -> %v (%T)", value, value, typedValue, typedValue)
			
			util.SetNestedValue(config.Data, configKey, typedValue)
			config.Sources[configKey] = "environment:" + key
		} else if strings.HasPrefix(key, e.varPrefix) {
			varName := strings.TrimPrefix(key, e.varPrefix)
			logger.Debug("  - Loading from env var: %s -> %s", key, varName)
			vars[varName] = value
		}
	}

	return &LoadResult{
		Config: config,
		Vars:   vars,
	}
}

// GetKeyPrefix returns the current key prefix.
func (e *Environment) GetKeyPrefix() string {
	return e.keyPrefix
}

// GetVarPrefix returns the current variable prefix.
func (e *Environment) GetVarPrefix() string {
	return e.varPrefix
}
