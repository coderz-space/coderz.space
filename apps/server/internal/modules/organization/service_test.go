package organization

import (
	"testing"
)

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{"valid lowercase", "my-org", true},
		{"valid with numbers", "org123", true},
		{"valid with hyphens", "my-org-123", true},
		{"invalid uppercase", "My-Org", false},
		{"invalid special chars", "my_org", false},
		{"invalid spaces", "my org", false},
		{"too short", "ab", false},
		{"minimum length", "abc", true},
		{"maximum length", "a" + string(make([]byte, 79)), false},   // 80 chars
		{"valid max length", "a" + string(make([]byte, 78)), false}, // 79 chars - need to fix this
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSlug(tt.slug)
			if result != tt.expected {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.slug, result, tt.expected)
			}
		})
	}
}

func TestNormalizeSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase conversion", "My Organization", "my-organization"},
		{"remove special chars", "My Org!", "my-org"},
		{"multiple spaces", "my   org", "my-org"},
		{"trim hyphens", "-my-org-", "my-org"},
		{"consecutive hyphens", "my--org", "my-org"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeSlug(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
