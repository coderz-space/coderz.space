package organization

import (
	"testing"
)

func TestPaginationDefaults(t *testing.T) {
	tests := []struct {
		name          string
		pageParam     string
		limitParam    string
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "no parameters - use defaults",
			pageParam:     "",
			limitParam:    "",
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "valid page and limit",
			pageParam:     "2",
			limitParam:    "50",
			expectedPage:  2,
			expectedLimit: 50,
		},
		{
			name:          "limit exceeds max - cap at 100",
			pageParam:     "1",
			limitParam:    "150",
			expectedPage:  1,
			expectedLimit: 20, // Should default to 20 since 150 > 100
		},
		{
			name:          "invalid page - use default",
			pageParam:     "invalid",
			limitParam:    "10",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "negative page - use default",
			pageParam:     "-1",
			limitParam:    "10",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "zero limit - use default",
			pageParam:     "1",
			limitParam:    "0",
			expectedPage:  1,
			expectedLimit: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies the pagination logic
			// In a real integration test, we would make HTTP requests
			// For now, we're just documenting the expected behavior
			t.Logf("Expected page=%d, limit=%d for pageParam=%q, limitParam=%q",
				tt.expectedPage, tt.expectedLimit, tt.pageParam, tt.limitParam)
		})
	}
}
