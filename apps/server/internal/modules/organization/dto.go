package organization

import "github.com/jackc/pgx/v5/pgtype"

// Organization DTOs

type CreateOrganizationRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=120"`
	Slug        string `json:"slug" validate:"required,min=3,max=80,lowercase,alphanum_hyphen"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=120"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=80,lowercase,alphanum_hyphen"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type OrganizationData struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type OrganizationResponse struct {
	Success bool             `json:"success"`
	Data    OrganizationData `json:"data"`
}

type OrganizationListResponse struct {
	Success bool               `json:"success"`
	Data    []OrganizationData `json:"data"`
	Meta    *PaginationMeta    `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Organization Member DTOs

type AddMemberRequest struct {
	UserID string `json:"userId" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=admin mentor mentee"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin mentor mentee"`
}

type MemberData struct {
	ID             pgtype.UUID `json:"id"`
	OrganizationID pgtype.UUID `json:"organizationId"`
	UserID         pgtype.UUID `json:"userId"`
	Role           string      `json:"role"`
	JoinedAt       string      `json:"joinedAt"`
	Name           string      `json:"name,omitempty"`
	Email          string      `json:"email,omitempty"`
	AvatarUrl      string      `json:"avatarUrl,omitempty"`
}

type MemberResponse struct {
	Success bool       `json:"success"`
	Data    MemberData `json:"data"`
}

type MemberListResponse struct {
	Success bool            `json:"success"`
	Data    []MemberData    `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Data    map[string]any `json:"data"`
}
