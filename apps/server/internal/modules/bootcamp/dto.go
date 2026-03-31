package bootcamp

import "github.com/jackc/pgx/v5/pgtype"

// Bootcamp DTOs

type CreateBootcampRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=120"`
	Description string `json:"description" validate:"omitempty,max=500"`
	StartDate   string `json:"startDate" validate:"omitempty,datetime=2006-01-02"`
	EndDate     string `json:"endDate" validate:"omitempty,datetime=2006-01-02"`
	IsActive    *bool  `json:"isActive" validate:"omitempty"`
}

type UpdateBootcampRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=120"`
	Description string `json:"description" validate:"omitempty,max=500"`
	StartDate   string `json:"startDate" validate:"omitempty,datetime=2006-01-02"`
	EndDate     string `json:"endDate" validate:"omitempty,datetime=2006-01-02"`
	IsActive    *bool  `json:"isActive" validate:"omitempty"`
}

type BootcampData struct {
	ID             pgtype.UUID `json:"id"`
	OrganizationID pgtype.UUID `json:"organizationId"`
	CreatedBy      pgtype.UUID `json:"createdBy"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	StartDate      string      `json:"startDate,omitempty"`
	EndDate        string      `json:"endDate,omitempty"`
	IsActive       bool        `json:"isActive"`
	CreatedAt      string      `json:"createdAt"`
	UpdatedAt      string      `json:"updatedAt"`
}

type BootcampResponse struct {
	Success bool         `json:"success"`
	Data    BootcampData `json:"data"`
}

type BootcampListResponse struct {
	Success bool           `json:"success"`
	Data    []BootcampData `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Bootcamp Enrollment DTOs

type EnrollMemberRequest struct {
	OrganizationMemberID string `json:"organizationMemberId" validate:"required,uuid"`
	Role                 string `json:"role" validate:"required,oneof=mentor mentee"`
}

type UpdateEnrollmentRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=mentor mentee"`
}

type EnrollmentData struct {
	ID                   pgtype.UUID `json:"id"`
	BootcampID           pgtype.UUID `json:"bootcampId"`
	OrganizationMemberID pgtype.UUID `json:"organizationMemberId"`
	Role                 string      `json:"role"`
	Status               string      `json:"status"`
	EnrolledAt           string      `json:"enrolledAt"`
	Name                 string      `json:"name,omitempty"`
	Email                string      `json:"email,omitempty"`
	AvatarUrl            string      `json:"avatarUrl,omitempty"`
	OrgRole              string      `json:"orgRole,omitempty"`
}

type EnrollmentResponse struct {
	Success bool           `json:"success"`
	Data    EnrollmentData `json:"data"`
}

type EnrollmentListResponse struct {
	Success bool             `json:"success"`
	Data    []EnrollmentData `json:"data"`
	Meta    *PaginationMeta  `json:"meta,omitempty"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Data    map[string]any `json:"data"`
}
