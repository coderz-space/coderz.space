package problem

import (
	"regexp"
	"strings"
)

// NormalizeTagName normalizes tag names to lowercase with hyphens
// Examples:
//   - "Arrays" -> "arrays"
//   - "Dynamic Programming" -> "dynamic-programming"
//   - "Two Pointers" -> "two-pointers"
func NormalizeTagName(name string) string {
	// Convert to lowercase
	normalized := strings.ToLower(name)

	// Replace spaces with hyphens
	normalized = strings.ReplaceAll(normalized, " ", "-")

	// Remove any characters that are not alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	normalized = reg.ReplaceAllString(normalized, "")

	// Replace multiple consecutive hyphens with a single hyphen
	reg = regexp.MustCompile(`-+`)
	normalized = reg.ReplaceAllString(normalized, "-")

	// Trim hyphens from start and end
	normalized = strings.Trim(normalized, "-")

	return normalized
}
