package variables

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"os"
	"strings"
)

// DefaultResolver is the default implementation of the Resolver interface.
type DefaultResolver struct {
	vars map[string]string
}

// NewResolver creates a new variable resolver, processing sources in the correct order of precedence.
func NewResolver(envVars map[string]string, varsFromFile map[string]interface{}, schemaVars []Definition, config map[string]interface{}) (*DefaultResolver, error) {
	resolved := make(map[string]string)

	// 1. Highest precedence: Variables from KONFIGO_VAR_ environment variables.
	if envVars != nil {
		logger.Debug("Loading variables from environment (KONFIGO_VAR_...) (highest priority)")
		for k, v := range envVars {
			resolved[k] = v
		}
	}

	// 2. Second highest precedence: Variables from the -V file.
	if varsFromFile != nil {
		logger.Debug("Loading variables from --vars-file")
		for k, v := range varsFromFile {
			if _, exists := resolved[k]; exists {
				logger.Debug("  - Skipping var '%s' from file (already defined by environment)", k)
				continue
			}
			resolved[k] = fmt.Sprintf("%v", v)
		}
	}

	// 3. Lowest precedence: Variables defined in the schema.
	logger.Debug("Resolving variables from schema `vars` section")
	for _, varDef := range schemaVars {
		if _, exists := resolved[varDef.Name]; exists {
			logger.Debug("  - Skipping var '%s' from schema (already defined with higher precedence)", varDef.Name)
			continue
		}

		var val string
		var found bool

		if varDef.FromEnv != "" {
			val, found = os.LookupEnv(varDef.FromEnv)
		} else if varDef.FromPath != "" {
			if v, ok := util.GetNestedValue(config, varDef.FromPath); ok {
				val = fmt.Sprintf("%v", v)
				found = true
			}
		} else if varDef.Value != "" {
			val = varDef.Value
			found = true
		}

		if !found {
			if varDef.DefaultValue != "" {
				val = varDef.DefaultValue
			} else {
				return nil, fmt.Errorf("variable '%s' could not be resolved and has no default value", varDef.Name)
			}
		}
		resolved[varDef.Name] = val
	}

	return &DefaultResolver{vars: resolved}, nil
}

// SubstituteString performs ${VAR} replacement on a single string.
func (r *DefaultResolver) SubstituteString(input string) string {
	return VarRegex.ReplaceAllStringFunc(input, func(match string) string {
		varName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		if val, ok := r.vars[varName]; ok {
			return val
		}
		// Leave unresolved variables as-is
		return match
	})
}
