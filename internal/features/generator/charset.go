package generator

import mrand "math/rand"

// Character sets shared by the random and id generators.
const (
	// charsetAlphanumeric contains [a-zA-Z0-9].
	charsetAlphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// charsetAlpha contains [a-zA-Z].
	charsetAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// charsetNumeric contains [0-9].
	charsetNumeric = "0123456789"
)

// randomStringFromCharset builds a string of the given length by drawing
// uniformly from charset. Callers are responsible for validating length;
// a non-positive length yields an empty string.
func randomStringFromCharset(rng *mrand.Rand, charset string, length int) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}
