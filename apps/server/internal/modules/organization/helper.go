package organization

import (
	"regexp"
	"strings"
)

// ValidateSlug checks if a slug is valid (lowercase, alphanumeric with hyphens)
func ValidateSlug(slug string) bool {
	// Slug should be lowercase, alphanumeric with hyphens
	match, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return match && len(slug) >= 3 && len(slug) <= 80
}

// NormalizeSlug converts a string to a valid slug format
func NormalizeSlug(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}
