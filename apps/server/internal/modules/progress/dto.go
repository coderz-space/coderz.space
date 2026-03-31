package progress

import "github.com/jackc/pgx/v5/pgtype"

// Doubt DTOs

// CreateDoubtRequest represents the request body for creating a doubt
// @Description Request body for creating a doubt on an assignment problem
type CreateDoubtRequest struct {
	AssignmentProblemID string `json:"assignmentProblemId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message             string `json:"message" validate:"required,min=10,max=2000" example:"I'm having trouble understanding the time complexity of this algorithm"`
}

// ResolveDoubtRequest represents the request body for resolving a doubt
// @Description Request body for resolving a doubt with optional resolution note
type ResolveDoubtRequest struct {
	ResolutionNote string `json:"resolutionNote" validate:"omitempty,max=1000" example:"The time complexity is O(n log n) because of the sorting step"`
}

// DoubtData represents the doubt response data
// @Description Doubt details with resolution information
type DoubtData struct {
	ResolvedAt          string      `json:"resolvedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	Message             string      `json:"message" example:"I'm having trouble understanding the time complexity"`
	ResolutionNote      string      `json:"resolutionNote,omitempty" example:"The time complexity is O(n log n)"`
	CreatedAt           string      `json:"createdAt" example:"2024-01-15T09:00:00Z"`
	UpdatedAt           string      `json:"updatedAt" example:"2024-01-15T10:30:00Z"`
	RaisedByName        string      `json:"raisedByName,omitempty" example:"John Doe"`
	RaisedByEmail       string      `json:"raisedByEmail,omitempty" example:"john@example.com"`
	ResolvedByName      string      `json:"resolvedByName,omitempty" example:"Jane Smith"`
	AssignmentProblemID pgtype.UUID `json:"assignmentProblemId" example:"550e8400-e29b-41d4-a716-446655440001"`
	RaisedBy            pgtype.UUID `json:"raisedBy" example:"550e8400-e29b-41d4-a716-446655440002"`
	ResolvedBy          pgtype.UUID `json:"resolvedBy,omitempty" example:"550e8400-e29b-41d4-a716-446655440003"`
	ID                  pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Resolved            bool        `json:"resolved" example:"false"`
}

// DoubtResponse represents a single doubt response
// @Description Response containing a single doubt
type DoubtResponse struct {
	Data    DoubtData `json:"data"`
	Success bool      `json:"success" example:"true"`
}

// DoubtListResponse represents a list of doubts with pagination
// @Description Response containing a list of doubts with cursor-based pagination
type DoubtListResponse struct {
	Meta    *CursorPagination `json:"meta,omitempty"`
	Data    []DoubtData       `json:"data"`
	Success bool              `json:"success" example:"true"`
}

// CursorPagination represents cursor-based pagination metadata
// @Description Cursor-based pagination metadata for large datasets
type CursorPagination struct {
	NextCursor string `json:"nextCursor,omitempty" example:"eyJpZCI6IjU1MGU4NDAwLWUyOWItNDFkNC1hNzE2LTQ0NjY1NTQ0MDAwMCJ9"`
	HasMore    bool   `json:"hasMore" example:"true"`
	Limit      int    `json:"limit" example:"20"`
}

// GenericResponse represents a generic success response
// @Description Generic success response
type GenericResponse struct {
	Data    map[string]any `json:"data"`
	Success bool           `json:"success" example:"true"`
}
