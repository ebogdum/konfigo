package generator

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	"math/rand"
	"strconv"
	"time"
)

// RandomGeneratorType is the type identifier for the random generator.
const RandomGeneratorType = "random"

// RandomGenerator generates random values.
type RandomGenerator struct{}

// Type returns the generator type.
func (g *RandomGenerator) Type() string {
	return RandomGeneratorType
}

// Generate implements the random generator logic.
// It generates random values based on the format specification.
// Supported formats:
// - "int:min:max": Random integer between min and max (inclusive)
// - "float:min:max": Random float between min and max
// - "string:length": Random string of specified length using [a-zA-Z0-9]
// - "bytes:length": Random bytes as hex string
// - "uuid": UUID v4 format (8-4-4-4-12 hex digits)
func (g *RandomGenerator) Generate(config map[string]interface{}, def Definition, resolver VariableResolver) error {
	logger.Debug("  - Applying random generator for target path '%s'", def.TargetPath)

	format := def.Format
	if format == "" {
		return fmt.Errorf("random generator: format is required")
	}

	// Initialize random seed based on current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var result string
	var err error

	switch {
	case format == "uuid":
		result = generateUUID(rng)
	case len(format) >= 4 && format[:4] == "int:":
		result, err = generateRandomInt(rng, format[4:])
	case len(format) >= 6 && format[:6] == "float:":
		result, err = generateRandomFloat(rng, format[6:])
	case len(format) >= 7 && format[:7] == "string:":
		result, err = generateRandomString(rng, format[7:])
	case len(format) >= 6 && format[:6] == "bytes:":
		result, err = generateRandomBytes(rng, format[6:])
	default:
		return fmt.Errorf("random generator: unsupported format '%s'", format)
	}

	if err != nil {
		return fmt.Errorf("random generator: %w", err)
	}

	// Substitute any global variables in the final result
	if resolver != nil {
		result = resolver.SubstituteString(result)
	}

	// Set the generated value in the configuration
	util.SetNestedValue(config, def.TargetPath, result)

	logger.Debug("    Generated random value at path '%s'", def.TargetPath)
	return nil
}

// generateUUID creates a UUID v4 format string
func generateUUID(rng *rand.Rand) string {
	// Generate 16 random bytes
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(rng.Intn(256))
	}

	// Set version (4) and variant bits according to RFC 4122
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// generateRandomInt creates a random integer in the specified range
func generateRandomInt(rng *rand.Rand, params string) (string, error) {
	parts := splitParams(params, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("int format requires 'min:max', got '%s'", params)
	}

	min, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid min value '%s': %w", parts[0], err)
	}

	max, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid max value '%s': %w", parts[1], err)
	}

	if min > max {
		return "", fmt.Errorf("min (%d) cannot be greater than max (%d)", min, max)
	}

	result := rng.Intn(max-min+1) + min
	return strconv.Itoa(result), nil
}

// generateRandomFloat creates a random float in the specified range
func generateRandomFloat(rng *rand.Rand, params string) (string, error) {
	parts := splitParams(params, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("float format requires 'min:max', got '%s'", params)
	}

	min, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", fmt.Errorf("invalid min value '%s': %w", parts[0], err)
	}

	max, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", fmt.Errorf("invalid max value '%s': %w", parts[1], err)
	}

	if min > max {
		return "", fmt.Errorf("min (%f) cannot be greater than max (%f)", min, max)
	}

	result := rng.Float64()*(max-min) + min
	return fmt.Sprintf("%.6f", result), nil
}

// generateRandomString creates a random string of specified length using [a-zA-Z0-9]
func generateRandomString(rng *rand.Rand, params string) (string, error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", params, err)
	}

	if length < 0 {
		return "", fmt.Errorf("length cannot be negative: %d", length)
	}

	if length == 0 {
		return "", nil
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result), nil
}

// generateRandomBytes creates random bytes as hex string
func generateRandomBytes(rng *rand.Rand, params string) (string, error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", params, err)
	}

	if length < 0 {
		return "", fmt.Errorf("length cannot be negative: %d", length)
	}

	if length == 0 {
		return "", nil
	}

	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = byte(rng.Intn(256))
	}

	return fmt.Sprintf("%x", bytes), nil
}

// splitParams splits parameter string by delimiter
func splitParams(s, delim string) []string {
	if s == "" {
		return []string{}
	}
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+len(delim)] == delim {
			parts = append(parts, s[start:i])
			start = i + len(delim)
			i += len(delim) - 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// ValidateDefinition validates a random generator definition.
func (g *RandomGenerator) ValidateDefinition(def Definition) error {
	if def.TargetPath == "" {
		return fmt.Errorf("random generator: targetPath is required and cannot be empty")
	}

	if def.Format == "" {
		return fmt.Errorf("random generator: format is required and cannot be empty")
	}

	// Note: Sources should be empty for random generator
	if len(def.Sources) > 0 {
		return fmt.Errorf("random generator: sources should not be specified (random values are generated, not derived from config)")
	}

	return nil
}
