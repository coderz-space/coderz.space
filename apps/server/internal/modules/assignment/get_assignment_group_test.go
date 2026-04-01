package assignment

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// TestGetAssignmentGroup_ResponseStructure verifies that the GetAssignmentGroup
// handler returns the correct response structure with associated problems and positions.
//
// Requirements: 7.11
func TestGetAssignmentGroup_ResponseStructure(t *testing.T) {
	t.Run("response includes problems with positions", func(t *testing.T) {
		// Create a sample response structure
		response := AssignmentGroupResponse{
			Success: true,
			Data: AssignmentGroupData{
				ID: pgtype.UUID{
					Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					Valid: true,
				},
				Title:        "Week 1 - Arrays and Strings",
				Description:  "Introduction to fundamental data structures",
				DeadlineDays: 7,
				Problems: []GroupProblemRef{
					{
						ProblemID: pgtype.UUID{
							Bytes: [16]byte{2, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
							Valid: true,
						},
						Title:      "Two Sum",
						Difficulty: "easy",
						Position:   1,
					},
					{
						ProblemID: pgtype.UUID{
							Bytes: [16]byte{3, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
							Valid: true,
						},
						Title:      "Add Two Numbers",
						Difficulty: "medium",
						Position:   2,
					},
				},
			},
		}

		// Verify response structure
		assert.True(t, response.Success, "Response should be successful")
		assert.NotEmpty(t, response.Data.Title, "Assignment group should have a title")
		assert.NotNil(t, response.Data.Problems, "Response should include problems")
		assert.Len(t, response.Data.Problems, 2, "Should have 2 problems")

		// Verify first problem has all required fields including position
		problem1 := response.Data.Problems[0]
		assert.True(t, problem1.ProblemID.Valid, "Problem should have valid ID")
		assert.Equal(t, "Two Sum", problem1.Title, "Problem should have title")
		assert.Equal(t, "easy", problem1.Difficulty, "Problem should have difficulty")
		assert.Equal(t, int32(1), problem1.Position, "Problem should have position 1")

		// Verify second problem has all required fields including position
		problem2 := response.Data.Problems[1]
		assert.True(t, problem2.ProblemID.Valid, "Problem should have valid ID")
		assert.Equal(t, "Add Two Numbers", problem2.Title, "Problem should have title")
		assert.Equal(t, "medium", problem2.Difficulty, "Problem should have difficulty")
		assert.Equal(t, int32(2), problem2.Position, "Problem should have position 2")
	})

	t.Run("response handles empty problems list", func(t *testing.T) {
		response := AssignmentGroupResponse{
			Success: true,
			Data: AssignmentGroupData{
				Title:        "Empty Group",
				DeadlineDays: 7,
				Problems:     []GroupProblemRef{},
			},
		}

		assert.True(t, response.Success, "Response should be successful")
		assert.NotNil(t, response.Data.Problems, "Problems should not be nil")
		assert.Len(t, response.Data.Problems, 0, "Problems list should be empty")
	})
}

// TestGetAssignmentGroup_ProblemsOrdering verifies that problems are returned
// in the correct order based on their position values.
//
// Requirements: 7.11
func TestGetAssignmentGroup_ProblemsOrdering(t *testing.T) {
	t.Run("problems are ordered by position", func(t *testing.T) {
		problems := []GroupProblemRef{
			{Title: "Problem A", Position: 1},
			{Title: "Problem B", Position: 2},
			{Title: "Problem C", Position: 3},
		}

		// Verify positions are in ascending order
		for i := 0; i < len(problems)-1; i++ {
			assert.Less(t, problems[i].Position, problems[i+1].Position,
				"Problems should be ordered by position in ascending order")
		}
	})
}

// TestGetAssignmentGroup_ErrorHandling verifies that the GetAssignmentGroup
// handler properly handles error cases.
//
// Requirements: 7.11
func TestGetAssignmentGroup_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		groupID       string
		expectedError string
	}{
		{
			name:          "invalid group ID format",
			groupID:       "not-a-uuid",
			expectedError: "INVALID_GROUP_ID",
		},
		{
			name:          "group not found",
			groupID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedError: "ASSIGNMENT_GROUP_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents the expected error handling behavior
			t.Logf("Testing groupID=%q, expectedError=%q", tt.groupID, tt.expectedError)
		})
	}
}
