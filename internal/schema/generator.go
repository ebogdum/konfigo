package schema

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"strings"
)

func ApplyGenerators(config map[string]interface{}, generators []GeneratorDef, resolver *Resolver) error {
	logger.Debug("Applying %d generators...", len(generators))
	for _, g := range generators {
		logger.Debug("  - Applying generator type '%s'", g.Type)
		switch g.Type {
		case "concat":
			var replacerArgs []string
			for placeholder, path := range g.Sources {
				val, found := util.GetNestedValue(config, path)
				if !found {
					return fmt.Errorf("concat generator: source path '%s' not found", path)
				}
				replacerArgs = append(replacerArgs, fmt.Sprintf("{%s}", placeholder), fmt.Sprintf("%v", val))
			}
			replacer := strings.NewReplacer(replacerArgs...)
			result := replacer.Replace(g.Format)
			// Substitute global vars in the final result
			result = SubstituteString(result, resolver)
			util.SetNestedValue(config, g.TargetPath, result)
		default:
			return fmt.Errorf("unsupported generator type: %s", g.Type)
		}
	}
	return nil
}
