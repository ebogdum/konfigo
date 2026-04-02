package generator

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	mrand "math/rand"
	"strconv"
)

// RandomGeneratorType is the type identifier for the random generator.
const RandomGeneratorType = "random"

// RandomGenerator generates random values.
type RandomGenerator struct{}

// Type returns the generator type.
func (g *RandomGenerator) Type() string {
	return RandomGeneratorType
}

// newCryptoRand returns a math/rand.Rand seeded from crypto/rand for
// cryptographically-seeded pseudo-random generation.
func newCryptoRand() (*mrand.Rand, error) {
	var seed int64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &seed); err != nil {
		return nil, fmt.Errorf("failed to seed random generator from crypto/rand: %w", err)
	}
	return mrand.New(mrand.NewSource(seed)), nil
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

	rng, err := newCryptoRand()
	if err != nil {
		return fmt.Errorf("random generator: %w", err)
	}

	var result string
	var genErr error

	switch {
	case format == "uuid":
		result = generateUUID()
	case len(format) >= 4 && format[:4] == "int:":
		result, genErr = generateRandomInt(rng, format[4:])
	case len(format) >= 6 && format[:6] == "float:":
		result, genErr = generateRandomFloat(rng, format[6:])
	case len(format) >= 7 && format[:7] == "string:":
		result, genErr = generateRandomString(rng, format[7:])
	case len(format) >= 6 && format[:6] == "bytes:":
		result, genErr = generateRandomBytes(format[6:])
	default:
		return fmt.Errorf("random generator: unsupported format '%s'", format)
	}

	if genErr != nil {
		return fmt.Errorf("random generator: %w", genErr)
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

// generateUUID creates a UUID v4 format string using crypto/rand
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	// Set version (4) and variant bits according to RFC 4122
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// generateRandomInt creates a random integer in the specified range
func generateRandomInt(rng *mrand.Rand, params string) (string, error) {
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

	rangeSize := max - min + 1
	if rangeSize <= 0 {
		return "", fmt.Errorf("range [%d, %d] is too large or overflows", min, max)
	}

	result := rng.Intn(rangeSize) + min
	return strconv.Itoa(result), nil
}

// generateRandomFloat creates a random float in the specified range
func generateRandomFloat(rng *mrand.Rand, params string) (string, error) {
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
func generateRandomString(rng *mrand.Rand, params string) (string, error) {
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

// generateRandomBytes creates random bytes as hex string using crypto/rand
func generateRandomBytes(params string) (string, error) {
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

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return fmt.Sprintf("%x", b), nil
}

// splitParams splits parameter string by delimiter
func splitParams(s, delim string) []string {
	if s == "" {
		return []string{}
	}
	var parts []string
	start := 0
	for i := 0; i <= len(s)-len(delim); i++ {
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
