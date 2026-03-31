package problem

import "github.com/jackc/pgx/v5/pgtype"

// Problem DTOs

type CreateProblemRequest struct {
	Title        string `json:"title" validate:"required,min=3,max=200" example:"Two Sum"`
	Description  string `json:"description" validate:"required,min=10" example:"Given an array of integers, return indices of the two numbers that add up to a specific target."`
	Difficulty   string `json:"difficulty" validate:"required,oneof=easy medium hard" example:"easy"`
	ExternalLink string `json:"externalLink" validate:"omitempty,url" example:"https://leetcode.com/problems/two-sum/"`
}

type UpdateProblemRequest struct {
	Title        string `json:"title" validate:"omitempty,min=3,max=200" example:"Two Sum Updated"`
	Description  string `json:"description" validate:"omitempty,min=10" example:"Updated description"`
	Difficulty   string `json:"difficulty" validate:"omitempty,oneof=easy medium hard" example:"medium"`
	ExternalLink string `json:"externalLink" validate:"omitempty,url" example:"https://leetcode.com/problems/two-sum/"`
}

type ProblemData struct {
	ID             pgtype.UUID    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID pgtype.UUID    `json:"organizationId" example:"660e8400-e29b-41d4-a716-446655440000"`
	CreatedBy      pgtype.UUID    `json:"createdBy" example:"770e8400-e29b-41d4-a716-446655440000"`
	Title          string         `json:"title" example:"Two Sum"`
	Description    string         `json:"description" example:"Given an array of integers, return indices of the two numbers that add up to a specific target."`
	Difficulty     string         `json:"difficulty" example:"easy"`
	ExternalLink   string         `json:"externalLink,omitempty" example:"https://leetcode.com/problems/two-sum/"`
	CreatedAt      string         `json:"createdAt" example:"2024-01-01T10:00:00Z"`
	UpdatedAt      string         `json:"updatedAt" example:"2024-01-01T10:00:00Z"`
	ArchivedAt     string         `json:"archivedAt,omitempty" example:""`
	Tags           []TagData      `json:"tags,omitempty"`
	Resources      []ResourceData `json:"resources,omitempty"`
}

type ProblemResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    ProblemData `json:"data"`
}

type ProblemListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []ProblemData   `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

// Tag DTOs

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=80" example:"arrays"`
}

type UpdateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=80" example:"dynamic-programming"`
}

type AttachTagsRequest struct {
	TagIDs []string `json:"tagIds" validate:"required,min=1,dive,uuid" example:"550e8400-e29b-41d4-a716-446655440000,660e8400-e29b-41d4-a716-446655440000"`
}

type TagData struct {
	ID             pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID pgtype.UUID `json:"organizationId" example:"660e8400-e29b-41d4-a716-446655440000"`
	Name           string      `json:"name" example:"arrays"`
	CreatedAt      string      `json:"createdAt" example:"2024-01-01T10:00:00Z"`
}

type TagResponse struct {
	Success bool    `json:"success" example:"true"`
	Data    TagData `json:"data"`
}

type TagListResponse struct {
	Success bool      `json:"success" example:"true"`
	Data    []TagData `json:"data"`
}

// Resource DTOs

type CreateResourceRequest struct {
	Title string `json:"title" validate:"required,min=2,max=150" example:"Two Sum Solution Explanation"`
	URL   string `json:"url" validate:"required,url" example:"https://www.youtube.com/watch?v=example"`
}

type UpdateResourceRequest struct {
	Title string `json:"title" validate:"omitempty,min=2,max=150" example:"Updated Resource Title"`
	URL   string `json:"url" validate:"omitempty,url" example:"https://www.youtube.com/watch?v=updated"`
}

type ResourceData struct {
	ID        pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProblemID pgtype.UUID `json:"problemId" example:"660e8400-e29b-41d4-a716-446655440000"`
	Title     string      `json:"title" example:"Two Sum Solution Explanation"`
	URL       string      `json:"url" example:"https://www.youtube.com/watch?v=example"`
	CreatedAt string      `json:"createdAt" example:"2024-01-01T10:00:00Z"`
}

type ResourceResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    ResourceData `json:"data"`
}

type ResourceListResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []ResourceData `json:"data"`
}

// Common DTOs

type PaginationMeta struct {
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
	Total int `json:"total" example:"100"`
}

type GenericResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    map[string]any `json:"data"`
}
