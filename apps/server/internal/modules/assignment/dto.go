package assignment

import "github.com/jackc/pgx/v5/pgtype"

// Assignment Group DTOs

type CreateAssignmentGroupRequest struct {
	Title        string `json:"title" validate:"required,min=3,max=150" example:"Week 1 - Arrays and Strings"`
	Description  string `json:"description" validate:"omitempty,max=1000" example:"Introduction to fundamental data structures"`
	DeadlineDays int32  `json:"deadlineDays" validate:"required,min=1" example:"7"`
}

type UpdateAssignmentGroupRequest struct {
	Title        string `json:"title" validate:"omitempty,min=3,max=150" example:"Week 1 - Arrays and Strings (Updated)"`
	Description  string `json:"description" validate:"omitempty,max=1000" example:"Updated description"`
	DeadlineDays int32  `json:"deadlineDays" validate:"omitempty,min=1" example:"10"`
}

type AddProblemsToGroupRequest struct {
	Problems []GroupProblemInput `json:"problems" validate:"required,min=1,dive"`
}

type ReplaceGroupProblemsRequest struct {
	Problems []GroupProblemInput `json:"problems" validate:"required,min=1,dive"`
}

type GroupProblemInput struct {
	ProblemID string `json:"problemId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Position  int32  `json:"position" validate:"required,min=1" example:"1"`
}

type AssignmentGroupData struct {
	ID           pgtype.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BootcampID   pgtype.UUID       `json:"bootcampId" example:"660e8400-e29b-41d4-a716-446655440000"`
	CreatedBy    pgtype.UUID       `json:"createdBy" example:"770e8400-e29b-41d4-a716-446655440000"`
	Title        string            `json:"title" example:"Week 1 - Arrays and Strings"`
	Description  string            `json:"description,omitempty" example:"Introduction to fundamental data structures"`
	DeadlineDays int32             `json:"deadlineDays" example:"7"`
	CreatedAt    string            `json:"createdAt" example:"2024-01-01T10:00:00Z"`
	UpdatedAt    string            `json:"updatedAt" example:"2024-01-01T10:00:00Z"`
	Problems     []GroupProblemRef `json:"problems,omitempty"`
}

type GroupProblemRef struct {
	ProblemID  pgtype.UUID `json:"problemId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title      string      `json:"title" example:"Two Sum"`
	Difficulty string      `json:"difficulty" example:"easy"`
	Position   int32       `json:"position" example:"1"`
}

type AssignmentGroupResponse struct {
	Data    AssignmentGroupData `json:"data"`
	Success bool                `json:"success" example:"true"`
}

type AssignmentGroupListResponse struct {
	Meta    *PaginationMeta       `json:"meta,omitempty"`
	Data    []AssignmentGroupData `json:"data"`
	Success bool                  `json:"success" example:"true"`
}

// Assignment Instance DTOs

type CreateAssignmentRequest struct {
	AssignmentGroupID    string `json:"assignmentGroupId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	BootcampEnrollmentID string `json:"bootcampEnrollmentId" validate:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440000"`
	DeadlineAt           string `json:"deadlineAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" example:"2024-01-15T23:59:59Z"`
}

type UpdateAssignmentRequest struct {
	DeadlineAt string `json:"deadlineAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" example:"2024-01-20T23:59:59Z"`
	Status     string `json:"status" validate:"omitempty,oneof=active completed expired" example:"completed"`
}

type UpdateAssignmentDeadlineRequest struct {
	DeadlineAt string `json:"deadlineAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2024-01-20T23:59:59Z"`
}

type UpdateAssignmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active completed expired" example:"completed"`
}

type AssignmentData struct {
	ID                   pgtype.UUID             `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	AssignmentGroupID    pgtype.UUID             `json:"assignmentGroupId" example:"660e8400-e29b-41d4-a716-446655440000"`
	BootcampEnrollmentID pgtype.UUID             `json:"bootcampEnrollmentId" example:"770e8400-e29b-41d4-a716-446655440000"`
	AssignedBy           pgtype.UUID             `json:"assignedBy" example:"880e8400-e29b-41d4-a716-446655440000"`
	AssignedAt           string                  `json:"assignedAt" example:"2024-01-01T10:00:00Z"`
	DeadlineAt           string                  `json:"deadlineAt,omitempty" example:"2024-01-08T23:59:59Z"`
	Status               string                  `json:"status" example:"active"`
	CreatedAt            string                  `json:"createdAt" example:"2024-01-01T10:00:00Z"`
	UpdatedAt            string                  `json:"updatedAt" example:"2024-01-01T10:00:00Z"`
	GroupTitle           string                  `json:"groupTitle,omitempty" example:"Week 1 - Arrays and Strings"`
	Problems             []AssignmentProblemData `json:"problems,omitempty"`
}

type AssignmentResponse struct {
	Data    AssignmentData `json:"data"`
	Success bool           `json:"success" example:"true"`
}

type AssignmentListResponse struct {
	Meta    *PaginationMeta  `json:"meta,omitempty"`
	Data    []AssignmentData `json:"data"`
	Success bool             `json:"success" example:"true"`
}

// Assignment Problem Progress DTOs

type UpdateAssignmentProblemRequest struct {
	Status       string `json:"status" validate:"omitempty,oneof=pending attempted completed" example:"completed"`
	SolutionLink string `json:"solutionLink" validate:"omitempty,url" example:"https://github.com/user/solution"`
	Notes        string `json:"notes" validate:"omitempty,max=2000" example:"Used dynamic programming approach"`
}

type AssignmentProblemData struct {
	ID           pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	AssignmentID pgtype.UUID `json:"assignmentId" example:"660e8400-e29b-41d4-a716-446655440000"`
	ProblemID    pgtype.UUID `json:"problemId" example:"770e8400-e29b-41d4-a716-446655440000"`
	Status       string      `json:"status" example:"pending"`
	SolutionLink string      `json:"solutionLink,omitempty" example:"https://github.com/user/solution"`
	Notes        string      `json:"notes,omitempty" example:"Used dynamic programming approach"`
	CompletedAt  string      `json:"completedAt,omitempty" example:"2024-01-05T14:30:00Z"`
	CreatedAt    string      `json:"createdAt" example:"2024-01-01T10:00:00Z"`
	UpdatedAt    string      `json:"updatedAt" example:"2024-01-05T14:30:00Z"`
	Title        string      `json:"title,omitempty" example:"Two Sum"`
	Difficulty   string      `json:"difficulty,omitempty" example:"easy"`
}

type AssignmentProblemResponse struct {
	Data    AssignmentProblemData `json:"data"`
	Success bool                  `json:"success" example:"true"`
}

type AssignmentProblemListResponse struct {
	Data    []AssignmentProblemData `json:"data"`
	Success bool                    `json:"success" example:"true"`
}

// Common DTOs

type PaginationMeta struct {
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
	Total int `json:"total" example:"100"`
}

type GenericResponse struct {
	Data    map[string]any `json:"data"`
	Success bool           `json:"success" example:"true"`
}
