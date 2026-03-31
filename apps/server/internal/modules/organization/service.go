package organization

import (
	"context"
	"errors"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	config  *config.Config
	pool    *pgxpool.Pool
}

func NewService(queries *db.Queries, config *config.Config, pool *pgxpool.Pool) *Service {
	return &Service{
		queries: queries,
		config:  config,
		pool:    pool,
	}
}

// Organization operations

func (s *Service) CreateOrganization(ctx context.Context, req CreateOrganizationRequest, userID pgtype.UUID) (*OrganizationData, error) {
	// Validate slug format
	if !ValidateSlug(req.Slug) {
		return nil, errors.New("INVALID_SLUG_FORMAT")
	}

	// Check if slug already exists
	_, err := s.queries.GetOrganizationBySlug(ctx, req.Slug)
	if err == nil {
		return nil, errors.New("SLUG_ALREADY_EXISTS")
	}

	// Use transaction to ensure atomicity
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Create organization with PENDING_APPROVAL status
	org, err := qtx.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Status:      db.OrgStatusPendingApproval,
	})
	if err != nil {
		return nil, err
	}

	// Add creator as admin member
	_, err = qtx.AddOrganizationMember(ctx, db.AddOrganizationMemberParams{
		OrganizationID: org.ID,
		UserID:         userID,
		Role:           db.OrgMemberRoleAdmin,
	})
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.mapOrganizationToData(org), nil
}

func (s *Service) GetOrganizationByID(ctx context.Context, orgID pgtype.UUID) (*OrganizationData, error) {
	org, err := s.queries.GetOrganizationById(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return s.mapOrganizationToData(org), nil
}

func (s *Service) ListUserOrganizations(ctx context.Context, userID pgtype.UUID, page, limit int) ([]OrganizationData, int, error) {
	// Calculate offset from page and limit
	offset := (page - 1) * limit

	// Get total count
	count, err := s.queries.CountUserOrganizations(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated organizations
	orgs, err := s.queries.ListOrganizations(ctx, db.ListOrganizationsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]OrganizationData, len(orgs))
	for i, org := range orgs {
		result[i] = *s.mapOrganizationToData(org)
	}

	return result, int(count), nil
}

func (s *Service) UpdateOrganization(ctx context.Context, orgID pgtype.UUID, req UpdateOrganizationRequest) (*OrganizationData, error) {
	// Validate at least one field is provided
	if req.Name == "" && req.Slug == "" && req.Description == "" {
		return nil, errors.New("NO_FIELDS_PROVIDED")
	}

	// If slug is being updated, validate format and uniqueness
	if req.Slug != "" {
		if !ValidateSlug(req.Slug) {
			return nil, errors.New("INVALID_SLUG_FORMAT")
		}

		// Check if slug already exists (excluding current organization)
		existingOrg, err := s.queries.GetOrganizationBySlug(ctx, req.Slug)
		if err == nil && existingOrg.ID != orgID {
			return nil, errors.New("SLUG_ALREADY_EXISTS")
		}
	}

	org, err := s.queries.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:          orgID,
		Name:        pgtype.Text{String: req.Name, Valid: req.Name != ""},
		Slug:        pgtype.Text{String: req.Slug, Valid: req.Slug != ""},
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Status:      db.NullOrgStatus{Valid: false},
	})
	if err != nil {
		return nil, err
	}

	return s.mapOrganizationToData(org), nil
}

func (s *Service) ApproveOrganization(ctx context.Context, orgID pgtype.UUID) (*OrganizationData, error) {
	// First, get the organization to validate its current status
	existingOrg, err := s.queries.GetOrganizationById(ctx, orgID)
	if err != nil {
		return nil, errors.New("ORGANIZATION_NOT_FOUND")
	}

	// Validate organization is in PENDING_APPROVAL status
	if existingOrg.Status != db.OrgStatusPendingApproval {
		return nil, errors.New("ORGANIZATION_NOT_PENDING")
	}

	// Update status to APPROVED
	org, err := s.queries.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:          orgID,
		Name:        pgtype.Text{Valid: false},
		Description: pgtype.Text{Valid: false},
		Status:      db.NullOrgStatus{OrgStatus: db.OrgStatusApproved, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return s.mapOrganizationToData(org), nil
}

func (s *Service) GetPendingOrganizations(ctx context.Context) ([]OrganizationData, error) {
	orgs, err := s.queries.GetPendingOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]OrganizationData, len(orgs))
	for i, org := range orgs {
		result[i] = *s.mapOrganizationToData(org)
	}

	return result, nil
}

// Member operations

func (s *Service) AddMember(ctx context.Context, orgID pgtype.UUID, req AddMemberRequest) (*MemberData, error) {
	userUUID, err := utils.StringToUUID(req.UserID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID")
	}

	role, err := s.parseOrgMemberRole(req.Role)
	if err != nil {
		return nil, err
	}

	member, err := s.queries.AddOrganizationMember(ctx, db.AddOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userUUID,
		Role:           role,
	})
	if err != nil {
		return nil, err
	}

	return s.mapMemberToData(member), nil
}

func (s *Service) ListMembers(ctx context.Context, orgID pgtype.UUID) ([]MemberData, error) {
	members, err := s.queries.ListOrganizationMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	result := make([]MemberData, len(members))
	for i, member := range members {
		result[i] = MemberData{
			ID:             member.ID,
			OrganizationID: member.OrganizationID,
			UserID:         member.UserID,
			Role:           string(member.Role),
			JoinedAt:       member.JoinedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			Name:           member.Name,
			Email:          member.Email.String,
			AvatarUrl:      member.AvatarUrl.String,
		}
	}

	return result, nil
}

func (s *Service) UpdateMemberRole(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID, req UpdateMemberRoleRequest) (*MemberData, error) {
	role, err := s.parseOrgMemberRole(req.Role)
	if err != nil {
		return nil, err
	}

	member, err := s.queries.UpdateOrganizationMemberRole(ctx, db.UpdateOrganizationMemberRoleParams{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
	})
	if err != nil {
		return nil, err
	}

	return s.mapMemberToData(member), nil
}

func (s *Service) RemoveMember(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID) error {
	return s.queries.RemoveOrganizationMember(ctx, db.RemoveOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
}

func (s *Service) GetMember(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID) (*MemberData, error) {
	member, err := s.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}

	return s.mapMemberToData(member), nil
}

// Helper methods

func (s *Service) mapOrganizationToData(org db.Organization) *OrganizationData {
	return &OrganizationData{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description.String,
		Status:      string(org.Status),
		CreatedAt:   org.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   org.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) mapMemberToData(member db.OrganizationMember) *MemberData {
	return &MemberData{
		ID:             member.ID,
		OrganizationID: member.OrganizationID,
		UserID:         member.UserID,
		Role:           string(member.Role),
		JoinedAt:       member.JoinedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) parseOrgMemberRole(role string) (db.OrgMemberRole, error) {
	switch role {
	case "admin":
		return db.OrgMemberRoleAdmin, nil
	case "mentor":
		return db.OrgMemberRoleMentor, nil
	case "mentee":
		return db.OrgMemberRoleMentee, nil
	default:
		return "", errors.New("INVALID_ROLE")
	}
}
