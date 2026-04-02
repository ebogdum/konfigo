package variables

import (
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// Substitute performs ${VAR} replacement on the entire configuration map.
func Substitute(config map[string]interface{}, resolver Resolver) map[string]interface{} {
	if config == nil {
		return nil
	}
	logger.Debug("Performing variable substitution...")
	replacerFunc := func(s string) string {
		return resolver.SubstituteString(s)
	}

	result := util.WalkAndReplace(config, replacerFunc)
	if result == nil {
		return nil
	}
	return result.(map[string]interface{})
}
