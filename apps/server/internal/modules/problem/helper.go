package problem

import (
	"regexp"
	"strings"
)

// NormalizeTagName normalizes tag names to lowercase with hyphens
// Example: "Dynamic Programming" -> "dynamic-programming"
func NormalizeTagName(name string) string {
	// Convert to lowercase
	normalized := strings.ToLower(name)

	// Replace spaces with hyphens
	normalized = strings.ReplaceAll(normalized, " ", "-")

	// Remove any characters that are not alphanumeric or hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	normalized = reg.ReplaceAllString(normalized, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	normalized = reg.ReplaceAllString(normalized, "-")

	// Trim hyphens from start and end
	normalized = strings.Trim(normalized, "-")

	return normalized
}

// ValidateDifficulty checks if difficulty is one of: easy, medium, hard
func ValidateDifficulty(difficulty string) bool {
	switch difficulty {
	case "easy", "medium", "hard":
		return true
	default:
		return false
	}
}

// DeduplicateStrings removes duplicate strings from a slice
func DeduplicateStrings(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
