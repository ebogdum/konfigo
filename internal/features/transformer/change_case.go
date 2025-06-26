package transformer

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"

	"github.com/iancoleman/strcase"
)

// ChangeCaseType is the type identifier for the change case transformer.
const ChangeCaseType = "changeCase"

// ChangeCaseTransformer changes the case of string values in configuration.
type ChangeCaseTransformer struct{}

// Type returns the transformer type.
func (t *ChangeCaseTransformer) Type() string {
	return ChangeCaseType
}

// Transform implements the change case transformation logic.
// It changes the case of a string value at the specified path.
func (t *ChangeCaseTransformer) Transform(config map[string]interface{}, def Definition) error {
	logger.Debug("  - Applying changeCase transform at path '%s' to case '%s'", def.Path, def.Case)
	
	// Get the value from the specified path
	value, found := util.GetNestedValue(config, def.Path)
	if !found {
		return fmt.Errorf("changeCase: path '%s' not found", def.Path)
	}
	
	// Ensure the value is a string
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("changeCase: value at path '%s' is not a string (got %T)", def.Path, value)
	}
	
	// Apply the case transformation
	newValue, err := t.applyCase(strValue, def.Case)
	if err != nil {
		return fmt.Errorf("changeCase: %w", err)
	}
	
	// Set the transformed value
	util.SetNestedValue(config, def.Path, newValue)
	
	logger.Debug("    Changed case from '%s' to '%s'", strValue, newValue)
	return nil
}

// applyCase applies the specified case transformation to a string.
func (t *ChangeCaseTransformer) applyCase(input, caseType string) (string, error) {
	switch strings.ToLower(caseType) {
	case "upper":
		return strings.ToUpper(input), nil
	case "lower":
		return strings.ToLower(input), nil
	case "snake":
		return strcase.ToSnake(input), nil
	case "camel":
		return strcase.ToCamel(input), nil
	case "kebab":
		return strcase.ToKebab(input), nil
	case "pascal":
		return strcase.ToCamel(input), nil // Pascal is same as camel in this library
	default:
		return "", fmt.Errorf("unsupported case type '%s'. Supported: upper, lower, snake, camel, kebab, pascal", caseType)
	}
}

// ValidateDefinition validates a change case transformer definition.
func (t *ChangeCaseTransformer) ValidateDefinition(def Definition) error {
	if def.Path == "" {
		return fmt.Errorf("changeCase transformer: 'path' is required")
	}
	
	if def.Case == "" {
		return fmt.Errorf("changeCase transformer: 'case' is required")
	}
	
	// Validate case type
	_, err := t.applyCase("test", def.Case)
	if err != nil {
		return fmt.Errorf("changeCase transformer: %w", err)
	}
	
	return nil
}
