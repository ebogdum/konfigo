package generator

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"
)

// ConcatGeneratorType is the type identifier for the concat generator.
const ConcatGeneratorType = "concat"

// ConcatGenerator generates values by concatenating placeholders with values from configuration paths.
type ConcatGenerator struct{}

// Type returns the generator type.
func (g *ConcatGenerator) Type() string {
	return ConcatGeneratorType
}

// Generate implements the concat generator logic.
// It replaces placeholders in the format string with values from specified configuration paths.
func (g *ConcatGenerator) Generate(config map[string]interface{}, def Definition, resolver VariableResolver) error {
	logger.Debug("  - Applying concat generator for target path '%s'", def.TargetPath)

	// Build replacement arguments for string replacer
	var replacerArgs []string
	for placeholder, sourcePath := range def.Sources {
		value, found := util.GetNestedValue(config, sourcePath)
		if !found {
			return fmt.Errorf("concat generator: source path '%s' not found in configuration", sourcePath)
		}

		// Add placeholder and its replacement value to the replacer arguments
		placeholderKey := fmt.Sprintf("{%s}", placeholder)
		replacementValue := fmt.Sprintf("%v", value)

		replacerArgs = append(replacerArgs, placeholderKey, replacementValue)
	}

	// Create replacer and apply to format string
	replacer := strings.NewReplacer(replacerArgs...)
	result := replacer.Replace(def.Format)

	// Check for any remaining unresolved {PLACEHOLDER} patterns in the result.
	// strings.NewReplacer leaves unmatched patterns as-is, so any {name} still
	// present after replacement is a missing source definition.
	// Skip ${VAR} patterns — those are variable substitution placeholders, not concat sources.
	var unresolvedPlaceholders []string
	remaining := result
	for {
		openIdx := strings.Index(remaining, "{")
		if openIdx < 0 {
			break
		}
		closeIdx := strings.Index(remaining[openIdx:], "}")
		if closeIdx < 0 {
			break
		}
		placeholder := remaining[openIdx : openIdx+closeIdx+1]
		// Skip ${VAR} variable substitution patterns (preceded by $)
		if openIdx > 0 && remaining[openIdx-1] == '$' {
			remaining = remaining[openIdx+closeIdx+1:]
			continue
		}
		unresolvedPlaceholders = append(unresolvedPlaceholders, placeholder)
		remaining = remaining[openIdx+closeIdx+1:]
	}
	if len(unresolvedPlaceholders) > 0 {
		return fmt.Errorf("concat generator: unresolved placeholder(s) %s in format string — no matching source defined", strings.Join(unresolvedPlaceholders, ", "))
	}

	// Substitute any global variables in the final result
	if resolver != nil {
		result = resolver.SubstituteString(result)
	}

	// Set the generated value in the configuration
	util.SetNestedValue(config, def.TargetPath, result)

	logger.Debug("    Generated value '%s' at path '%s'", result, def.TargetPath)
	return nil
}

// ValidateDefinition validates a concat generator definition.
func (g *ConcatGenerator) ValidateDefinition(def Definition) error {
	if def.TargetPath == "" {
		return fmt.Errorf("concat generator: targetPath is required and cannot be empty")
	}

	// Note: Empty format is allowed (creates empty string)
	// Note: Empty sources is allowed (for static text generation)

	return nil
}
