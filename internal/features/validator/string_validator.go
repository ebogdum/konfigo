package validator

import (
	"fmt"
	"regexp"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

// regexCacheMu protects eviction of the regex cache.
var regexCacheMu sync.Mutex

// regexCache caches compiled regexes to avoid recompilation per validation call.
var regexCache sync.Map

// regexCacheCount tracks the number of entries to enforce a size limit.
var regexCacheCount int64

// maxRegexCacheSize is the maximum number of cached regex patterns.
const maxRegexCacheSize = 500

// maxRegexInputLen is the maximum string length that will be tested against a regex
// to prevent catastrophic backtracking (ReDoS).
const maxRegexInputLen = 1 << 20 // 1 MiB

// getCompiledRegex returns a compiled regex, using a cache to avoid recompilation.
// The cache is bounded to maxRegexCacheSize entries; when full, new patterns are
// compiled but not cached until the cache is cleared.
func getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	count := atomic.AddInt64(&regexCacheCount, 1)
	if count <= maxRegexCacheSize {
		regexCache.Store(pattern, compiled)
	} else {
		// Cache is full — evict under lock to avoid data race on reassignment
		regexCacheMu.Lock()
		if atomic.LoadInt64(&regexCacheCount) > maxRegexCacheSize {
			regexCache.Range(func(key, _ interface{}) bool {
				regexCache.Delete(key)
				return true
			})
			atomic.StoreInt64(&regexCacheCount, 1)
		}
		regexCacheMu.Unlock()
		regexCache.Store(pattern, compiled)
	}
	return compiled, nil
}

// StringValidator validates string constraints (minLength, enum, regex).
type StringValidator struct{}

// Validate performs string validation.
func (sv *StringValidator) Validate(value interface{}, path string, rule Rule) error {
	// Skip if no string constraints
	if rule.MinLength == nil && len(rule.Enum) == 0 && rule.Regex == "" {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		if rule.MinLength != nil {
			return fmt.Errorf("path '%s': minLength validation requires a string, got %T", path, value)
		}
		if len(rule.Enum) > 0 {
			return fmt.Errorf("path '%s': enum validation requires a string, got %T", path, value)
		}
		if rule.Regex != "" {
			return fmt.Errorf("path '%s': regex validation requires a string, got %T", path, value)
		}
		return nil
	}

	// MinLength validation using rune count for correct Unicode character counting
	if rule.MinLength != nil {
		runeCount := utf8.RuneCountInString(str)
		if runeCount < *rule.MinLength {
			return fmt.Errorf("path '%s': length %d is less than minimum length %d", path, runeCount, *rule.MinLength)
		}
	}

	// Enum validation
	if len(rule.Enum) > 0 {
		match := false
		for _, enumVal := range rule.Enum {
			if str == enumVal {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("path '%s': value '%s' is not in the allowed list %v", path, str, rule.Enum)
		}
	}

	// Regex validation with cached compilation and input length cap
	if rule.Regex != "" {
		if len(str) > maxRegexInputLen {
			return fmt.Errorf("path '%s': value length %d exceeds maximum for regex validation (%d)", path, len(str), maxRegexInputLen)
		}
		compiledRegex, err := getCompiledRegex(rule.Regex)
		if err != nil {
			return fmt.Errorf("path '%s': invalid regex pattern '%s': %w", path, rule.Regex, err)
		}
		if !compiledRegex.MatchString(str) {
			return fmt.Errorf("path '%s': value '%s' does not match regex pattern '%s'", path, str, rule.Regex)
		}
	}

	return nil
}
