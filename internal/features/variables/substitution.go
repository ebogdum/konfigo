package variables

import (
	"konfigo/internal/logger"
	"konfigo/internal/util"
)

// Substitute performs ${VAR} replacement on the entire configuration map.
func Substitute(config map[string]interface{}, resolver Resolver) map[string]interface{} {
	logger.Debug("Performing variable substitution...")
	replacerFunc := func(s string) string {
		return resolver.SubstituteString(s)
	}

	return util.WalkAndReplace(config, replacerFunc).(map[string]interface{})
}
