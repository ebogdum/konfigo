package schema

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"

	"github.com/iancoleman/strcase"
)

func ApplyTransforms(config map[string]interface{}, transforms []TransformDef, resolver *Resolver) error {
	logger.Debug("Applying %d transformations...", len(transforms))
	for _, t := range transforms {
		// First, substitute any variables in the transform definition itself.
		t.Path = SubstituteString(t.Path, resolver)
		t.From = SubstituteString(t.From, resolver)
		t.To = SubstituteString(t.To, resolver)
		if s, ok := t.Value.(string); ok {
			t.Value = SubstituteString(s, resolver)
		}

		logger.Debug("  - Applying transform type '%s'", t.Type)
		switch t.Type {
		case "renameKey":
			val, found := util.GetNestedValue(config, t.From)
			if !found {
				return fmt.Errorf("renameKey: source path '%s' not found", t.From)
			}
			util.SetNestedValue(config, t.To, val)
			util.DeleteNestedValue(config, t.From)
		case "changeCase":
			val, found := util.GetNestedValue(config, t.Path)
			if !found {
				return fmt.Errorf("changeCase: path '%s' not found", t.Path)
			}
			strVal, ok := val.(string)
			if !ok {
				return fmt.Errorf("changeCase: value at path '%s' is not a string", t.Path)
			}
			var newStr string
			switch strings.ToLower(t.Case) {
			case "upper":
				newStr = strings.ToUpper(strVal)
			case "lower":
				newStr = strings.ToLower(strVal)
			case "snake":
				newStr = strcase.ToSnake(strVal)
			case "camel":
				newStr = strcase.ToCamel(strVal)
			default:
				return fmt.Errorf("changeCase: unsupported case '%s'", t.Case)
			}
			util.SetNestedValue(config, t.Path, newStr)
		case "addKeyPrefix":
			val, found := util.GetNestedValue(config, t.Path)
			if !found {
				return fmt.Errorf("addKeyPrefix: path '%s' not found", t.Path)
			}
			mapVal, ok := val.(map[string]interface{})
			if !ok {
				return fmt.Errorf("addKeyPrefix: value at path '%s' is not a map", t.Path)
			}
			newMap := make(map[string]interface{})
			for k, v := range mapVal {
				newMap[t.Prefix+k] = v
			}
			util.SetNestedValue(config, t.Path, newMap)
		case "setValue":
			util.SetNestedValue(config, t.Path, t.Value)
		default:
			return fmt.Errorf("unsupported transform type: %s", t.Type)
		}
	}
	return nil
}

// SubstituteString is a helper to perform substitution on a single string.
func SubstituteString(s string, resolver *Resolver) string {
	return VarRegex.ReplaceAllStringFunc(s, func(match string) string {
		varName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		if val, ok := resolver.vars[varName]; ok {
			return val
		}
		return match
	})
}
