package modules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// SwaggerDoc represents the structure of swagger.json
type SwaggerDoc struct {
	Paths map[string]map[string]interface{} `json:"paths"`
}

// TestSwaggerDocumentationCompleteness verifies that all endpoints with godoc comments
// appear in the generated Swagger documentation across all 5 modules.
//
// **Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5**
//
// This is a bug condition exploration test that MUST FAIL on unfixed code.
// The test verifies that ~25 endpoints wrapped with core.WithBody() are missing
// from the Swagger documentation because swaggo cannot parse generic wrapper functions.
//
// Expected outcome: TEST FAILS - documents which endpoints are missing
func TestSwaggerDocumentationCompleteness(t *testing.T) {
	// Read and parse swagger.json
	// The test runs from the module directory, so we need to go up to the server root
	swaggerPath := filepath.Join("..", "..", "swagger", "swagger.json")
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("Failed to read swagger.json: %v", err)
	}

	var swagger SwaggerDoc
	if err := json.Unmarshal(data, &swagger); err != nil {
		t.Fatalf("Failed to parse swagger.json: %v", err)
	}

	// Define all expected endpoints across all 5 modules
	// These endpoints have godoc comments and should appear in Swagger
	expectedEndpoints := map[string][]string{
		// Auth module - 4 endpoints wrapped with core.WithBody()
		"/v1/auth/signup":          {"post"},
		"/v1/auth/login":           {"post"},
		"/v1/auth/forgot-password": {"post"},
		"/v1/auth/reset-password":  {"post"},

		// Organization module - 4 endpoints wrapped with core.WithBody()
		"/v1/organizations":                          {"post"},
		"/v1/organizations/{orgId}":                  {"patch"},
		"/v1/organizations/{orgId}/members":          {"post"},
		"/v1/organizations/{orgId}/members/{userId}": {"patch"},

		// Bootcamp module - 4 endpoints wrapped with core.WithBody()
		"/v1/organizations/{orgId}/bootcamps":                                         {"post"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}":                            {"patch"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments":                {"post"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments/{enrollmentId}": {"patch"},

		// Problem module - 7 endpoints wrapped with core.WithBody()
		"/v1/organizations/{orgId}/problems":                                    {"post"},
		"/v1/organizations/{orgId}/problems/{problemId}":                        {"patch"},
		"/v1/organizations/{orgId}/tags":                                        {"post"},
		"/v1/organizations/{orgId}/tags/{tagId}":                                {"patch"},
		"/v1/organizations/{orgId}/problems/{problemId}/tags":                   {"post"},
		"/v1/organizations/{orgId}/problems/{problemId}/resources":              {"post"},
		"/v1/organizations/{orgId}/problems/{problemId}/resources/{resourceId}": {"patch"},

		// Assignment module - 6 endpoints wrapped with core.WithBody()
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups":                               {"post"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId}":                     {"patch"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId}/problems":            {"post"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments":                                     {"post"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId}":                      {"patch"},
		"/v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId}/problems/{problemId}": {"patch"},
	}

	// Track missing endpoints
	var missingEndpoints []string
	var missingMethods []string

	// Verify each endpoint exists with the correct HTTP method
	for path, methods := range expectedEndpoints {
		pathData, pathExists := swagger.Paths[path]

		if !pathExists {
			missingEndpoints = append(missingEndpoints, path)
			for _, method := range methods {
				missingMethods = append(missingMethods, method+" "+path)
			}
			continue
		}

		// Check if the required HTTP methods exist
		for _, method := range methods {
			if _, methodExists := pathData[method]; !methodExists {
				missingMethods = append(missingMethods, method+" "+path)
			}
		}
	}

	// Report findings
	if len(missingEndpoints) > 0 || len(missingMethods) > 0 {
		t.Errorf("\n=== BUG CONDITION CONFIRMED ===\n")
		t.Errorf("Swagger documentation is missing %d endpoints across all 5 modules\n", len(missingMethods))
		t.Errorf("\nMissing endpoints by module:\n")

		// Group by module for better reporting
		authMissing := 0
		orgMissing := 0
		bootcampMissing := 0
		problemMissing := 0
		assignmentMissing := 0

		t.Errorf("\nAuth module (expected 4):\n")
		for _, endpoint := range missingMethods {
			if contains(endpoint, "/v1/auth/") {
				t.Errorf("  - %s\n", endpoint)
				authMissing++
			}
		}

		t.Errorf("\nOrganization module (expected 4):\n")
		for _, endpoint := range missingMethods {
			if contains(endpoint, "/v1/organizations") &&
				!contains(endpoint, "/bootcamps") &&
				!contains(endpoint, "/problems") &&
				!contains(endpoint, "/tags") {
				t.Errorf("  - %s\n", endpoint)
				orgMissing++
			}
		}

		t.Errorf("\nBootcamp module (expected 4):\n")
		for _, endpoint := range missingMethods {
			if contains(endpoint, "/bootcamps") && !contains(endpoint, "/assignment") {
				t.Errorf("  - %s\n", endpoint)
				bootcampMissing++
			}
		}

		t.Errorf("\nProblem module (expected 7):\n")
		for _, endpoint := range missingMethods {
			if (contains(endpoint, "/problems") || contains(endpoint, "/tags") || contains(endpoint, "/resources")) &&
				!contains(endpoint, "/assignment") {
				t.Errorf("  - %s\n", endpoint)
				problemMissing++
			}
		}

		t.Errorf("\nAssignment module (expected 6):\n")
		for _, endpoint := range missingMethods {
			if contains(endpoint, "/assignment") {
				t.Errorf("  - %s\n", endpoint)
				assignmentMissing++
			}
		}

		t.Errorf("\n=== SUMMARY ===\n")
		t.Errorf("Total missing: %d out of 25 expected endpoints\n", len(missingMethods))
		t.Errorf("  Auth: %d/4\n", authMissing)
		t.Errorf("  Organization: %d/4\n", orgMissing)
		t.Errorf("  Bootcamp: %d/4\n", bootcampMissing)
		t.Errorf("  Problem: %d/7\n", problemMissing)
		t.Errorf("  Assignment: %d/6\n", assignmentMissing)
		t.Errorf("\nRoot cause: swaggo cannot parse handlers wrapped with core.WithBody() generic function\n")
		t.Errorf("This test MUST FAIL on unfixed code - failure confirms the bug exists\n")
	}

	// This assertion will fail on unfixed code, documenting the bug
	if len(missingMethods) > 0 {
		t.Fatalf("\nBug confirmed: %d endpoints are missing from Swagger documentation", len(missingMethods))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
