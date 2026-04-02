package generator

import (
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/util"
	mrand "math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// IdGeneratorType is the type identifier for the id generator.
const IdGeneratorType = "id"

// IdGenerator generates various types of ID values using [a-zA-Z0-9] characters.
type IdGenerator struct{}

// Type returns the generator type.
func (g *IdGenerator) Type() string {
	return IdGeneratorType
}

// Generate implements the id generator logic.
// It generates ID values based on the format specification.
// Supported formats:
// - "simple:length": Simple random ID of specified length using [a-zA-Z0-9]
// - "prefix:prefix:length": ID with specified prefix followed by random characters
// - "numeric:length": Numeric ID using only digits [0-9]
// - "alpha:length": Alphabetic ID using only letters [a-zA-Z]
// - "sequential": Sequential counter-based ID (starts from 1)
// - "timestamp": Timestamp-based ID (unix timestamp + random suffix)
func (g *IdGenerator) Generate(config map[string]interface{}, def Definition, resolver VariableResolver) error {
	logger.Debug("  - Applying id generator for target path '%s'", def.TargetPath)

	format := def.Format
	if format == "" {
		format = "simple:8" // Default format
	}

	rng, err := newCryptoRand()
	if err != nil {
		return fmt.Errorf("id generator: %w", err)
	}

	var result string

	switch {
	case format == "sequential":
		result = generateSequentialId(def.TargetPath)
	case format == "timestamp":
		result = generateTimestampId(rng)
	case len(format) >= 7 && format[:7] == "simple:":
		result, err = generateSimpleId(rng, format[7:])
	case len(format) >= 7 && format[:7] == "prefix:":
		result, err = generatePrefixId(rng, format[7:])
	case len(format) >= 8 && format[:8] == "numeric:":
		result, err = generateNumericId(rng, format[8:])
	case len(format) >= 6 && format[:6] == "alpha:":
		result, err = generateAlphaId(rng, format[6:])
	default:
		return fmt.Errorf("id generator: unsupported format '%s'", format)
	}

	if err != nil {
		return fmt.Errorf("id generator: %w", err)
	}

	// Substitute any global variables in the final result
	if resolver != nil {
		result = resolver.SubstituteString(result)
	}

	// Set the generated value in the configuration
	util.SetNestedValue(config, def.TargetPath, result)

	logger.Debug("    Generated ID '%s' at path '%s'", result, def.TargetPath)
	return nil
}

// generateSimpleId creates a simple random ID using [a-zA-Z0-9]
func generateSimpleId(rng *mrand.Rand, params string) (string, error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", params, err)
	}

	if length <= 0 {
		return "", fmt.Errorf("length must be positive: %d", length)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result), nil
}

// generatePrefixId creates an ID with a prefix followed by random characters
func generatePrefixId(rng *mrand.Rand, params string) (string, error) {
	parts := strings.SplitN(params, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("prefix format requires 'prefix:length', got '%s'", params)
	}

	prefix := parts[0]
	length, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", parts[1], err)
	}

	if length <= 0 {
		return "", fmt.Errorf("length must be positive: %d", length)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	suffix := make([]byte, length)
	for i := range suffix {
		suffix[i] = charset[rng.Intn(len(charset))]
	}

	return prefix + string(suffix), nil
}

// generateNumericId creates a numeric ID using only digits [0-9]
func generateNumericId(rng *mrand.Rand, params string) (string, error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", params, err)
	}

	if length <= 0 {
		return "", fmt.Errorf("length must be positive: %d", length)
	}

	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result), nil
}

// generateAlphaId creates an alphabetic ID using only letters [a-zA-Z]
func generateAlphaId(rng *mrand.Rand, params string) (string, error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		return "", fmt.Errorf("invalid length '%s': %w", params, err)
	}

	if length <= 0 {
		return "", fmt.Errorf("length must be positive: %d", length)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result), nil
}

// sequentialCountersMu protects sequentialCounters from concurrent access.
var sequentialCountersMu sync.Mutex

// sequentialCounters tracks sequential ID counters outside of the user config
// to avoid polluting the configuration data.
var sequentialCounters = make(map[string]int)

// generateSequentialId creates a sequential counter-based ID
func generateSequentialId(targetPath string) string {
	sequentialCountersMu.Lock()
	sequentialCounters[targetPath]++
	id := sequentialCounters[targetPath]
	sequentialCountersMu.Unlock()
	return strconv.Itoa(id)
}

// generateTimestampId creates a timestamp-based ID
func generateTimestampId(rng *mrand.Rand) string {
	timestamp := time.Now().Unix()

	// Add a 4-character random suffix to ensure uniqueness
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	suffix := make([]byte, 4)
	for i := range suffix {
		suffix[i] = charset[rng.Intn(len(charset))]
	}

	return fmt.Sprintf("%d%s", timestamp, string(suffix))
}

// ValidateDefinition validates an id generator definition.
func (g *IdGenerator) ValidateDefinition(def Definition) error {
	if def.TargetPath == "" {
		return fmt.Errorf("id generator: targetPath is required and cannot be empty")
	}

	// Note: Format is optional, defaults to simple:8
	// Note: Sources should be empty for id generator
	if len(def.Sources) > 0 {
		return fmt.Errorf("id generator: sources should not be specified (IDs are generated, not derived from config)")
	}

	return nil
}
