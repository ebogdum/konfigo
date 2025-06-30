package generator

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"time"
)

// TimestampGeneratorType is the type identifier for the timestamp generator.
const TimestampGeneratorType = "timestamp"

// TimestampGenerator generates timestamp values in various formats.
type TimestampGenerator struct{}

// Type returns the generator type.
func (g *TimestampGenerator) Type() string {
	return TimestampGeneratorType
}

// Generate implements the timestamp generator logic.
// It generates a timestamp string in the specified format.
// Supported formats:
// - "unix": Unix timestamp (seconds since epoch)
// - "unixmilli": Unix timestamp in milliseconds
// - "rfc3339": RFC3339 format (e.g., "2006-01-02T15:04:05Z07:00")
// - "iso8601": ISO8601 format (e.g., "2006-01-02T15:04:05Z")
// - Custom Go time format string
func (g *TimestampGenerator) Generate(config map[string]interface{}, def Definition, resolver VariableResolver) error {
	logger.Debug("  - Applying timestamp generator for target path '%s'", def.TargetPath)

	format := def.Format
	if format == "" {
		format = "rfc3339" // Default format
	}

	now := time.Now()
	var result string

	switch format {
	case "unix":
		result = fmt.Sprintf("%d", now.Unix())
	case "unixmilli":
		result = fmt.Sprintf("%d", now.UnixMilli())
	case "rfc3339":
		result = now.Format(time.RFC3339)
	case "iso8601":
		result = now.UTC().Format("2006-01-02T15:04:05Z")
	default:
		// Treat as custom Go time format
		result = now.Format(format)
	}

	// Substitute any global variables in the final result
	if resolver != nil {
		result = resolver.SubstituteString(result)
	}

	// Set the generated value in the configuration
	util.SetNestedValue(config, def.TargetPath, result)

	logger.Debug("    Generated timestamp '%s' at path '%s'", result, def.TargetPath)
	return nil
}

// ValidateDefinition validates a timestamp generator definition.
func (g *TimestampGenerator) ValidateDefinition(def Definition) error {
	if def.TargetPath == "" {
		return fmt.Errorf("timestamp generator: targetPath is required and cannot be empty")
	}

	// Note: Format is optional, defaults to rfc3339
	// Note: Sources should be empty for timestamp generator
	if len(def.Sources) > 0 {
		return fmt.Errorf("timestamp generator: sources should not be specified (timestamps are generated, not derived from config)")
	}

	return nil
}
