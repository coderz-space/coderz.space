package bootcamp

import (
	"context"
	"errors"

	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
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

// Bootcamp operations

func (s *Service) CreateBootcamp(ctx context.Context, orgID pgtype.UUID, req CreateBootcampRequest, createdBy pgtype.UUID) (*BootcampData, error) {
	// Validate organization exists and is APPROVED
	org, err := s.queries.GetOrganizationById(ctx, orgID)
	if err != nil {
		return nil, errors.New("ORGANIZATION_NOT_FOUND")
	}

	if org.Status != db.OrgStatusApproved {
		return nil, errors.New("ORGANIZATION_NOT_APPROVED")
	}

	// Validate date range if both dates are provided
	if !ValidateDateRange(req.StartDate, req.EndDate) {
		return nil, errors.New("INVALID_DATE_RANGE")
	}

	// Parse dates
	var startDate, endDate pgtype.Date
	if req.StartDate != "" {
		parsedStart, err := ParseDate(req.StartDate)
		if err != nil {
			return nil, errors.New("INVALID_START_DATE")
		}
		startDate = pgtype.Date{Time: parsedStart, Valid: true}
	}

	if req.EndDate != "" {
		parsedEnd, err := ParseDate(req.EndDate)
		if err != nil {
			return nil, errors.New("INVALID_END_DATE")
		}
		endDate = pgtype.Date{Time: parsedEnd, Valid: true}
	}

	// Default is_active to true if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	bootcamp, err := s.queries.CreateBootcamp(ctx, db.CreateBootcampParams{
		OrganizationID: orgID,
		CreatedBy:      createdBy,
		Name:           req.Name,
		Description:    pgtype.Text{String: req.Description, Valid: req.Description != ""},
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       isActive,
	})
	if err != nil {
		return nil, err
	}

	return s.mapBootcampToData(bootcamp), nil
}

func (s *Service) GetBootcampByID(ctx context.Context, bootcampID pgtype.UUID) (*BootcampData, error) {
	bootcamp, err := s.queries.GetBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, err
	}

	return s.mapBootcampToData(bootcamp), nil
}

func (s *Service) ListBootcampsByOrg(ctx context.Context, orgID pgtype.UUID) ([]BootcampData, error) {
	bootcamps, err := s.queries.ListBootcampsByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	result := make([]BootcampData, len(bootcamps))
	for i, bootcamp := range bootcamps {
		result[i] = *s.mapBootcampToData(bootcamp)
	}

	return result, nil
}

func (s *Service) ListBootcampsWithFilters(ctx context.Context, orgID pgtype.UUID, memberID *pgtype.UUID, isActive *bool, page, limit int) ([]BootcampData, int, error) {
	offset := (page - 1) * limit

	var bootcamps []db.Bootcamp
	var count int64
	var err error

	if memberID != nil && memberID.Valid {
		// Mentee: List bootcamps where they are enrolled
		bootcamps, err = s.queries.ListBootcampsByEnrollment(ctx, db.ListBootcampsByEnrollmentParams{
			OrganizationMemberID: *memberID,
			IsActive:             pgtype.Bool{Bool: isActive != nil && *isActive, Valid: isActive != nil},
			Limit:                int32(limit),
			Offset:               int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}

		count, err = s.queries.CountBootcampsByEnrollment(ctx, db.CountBootcampsByEnrollmentParams{
			OrganizationMemberID: *memberID,
			IsActive:             pgtype.Bool{Bool: isActive != nil && *isActive, Valid: isActive != nil},
		})
		if err != nil {
			return nil, 0, err
		}
	} else {
		// Admin/Mentor: List all bootcamps in organization
		bootcamps, err = s.queries.ListBootcampsByOrgWithPagination(ctx, db.ListBootcampsByOrgWithPaginationParams{
			OrganizationID: orgID,
			IsActive:       pgtype.Bool{Bool: isActive != nil && *isActive, Valid: isActive != nil},
			Limit:          int32(limit),
			Offset:         int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}

		count, err = s.queries.CountBootcampsByOrg(ctx, db.CountBootcampsByOrgParams{
			OrganizationID: orgID,
			IsActive:       pgtype.Bool{Bool: isActive != nil && *isActive, Valid: isActive != nil},
		})
		if err != nil {
			return nil, 0, err
		}
	}

	result := make([]BootcampData, len(bootcamps))
	for i, bootcamp := range bootcamps {
		result[i] = *s.mapBootcampToData(bootcamp)
	}

	return result, int(count), nil
}

func (s *Service) UpdateBootcamp(ctx context.Context, bootcampID pgtype.UUID, req UpdateBootcampRequest) (*BootcampData, error) {
	// Validate at least one field is provided
	if req.Name == "" && req.Description == "" && req.StartDate == "" && req.EndDate == "" && req.IsActive == nil {
		return nil, errors.New("NO_FIELDS_PROVIDED")
	}

	// Validate date range if both dates are provided
	if !ValidateDateRange(req.StartDate, req.EndDate) {
		return nil, errors.New("INVALID_DATE_RANGE")
	}

	// Parse dates
	var startDate, endDate pgtype.Date
	if req.StartDate != "" {
		parsedStart, err := ParseDate(req.StartDate)
		if err != nil {
			return nil, errors.New("INVALID_START_DATE")
		}
		startDate = pgtype.Date{Time: parsedStart, Valid: true}
	}

	if req.EndDate != "" {
		parsedEnd, err := ParseDate(req.EndDate)
		if err != nil {
			return nil, errors.New("INVALID_END_DATE")
		}
		endDate = pgtype.Date{Time: parsedEnd, Valid: true}
	}

	bootcamp, err := s.queries.UpdateBootcamp(ctx, db.UpdateBootcampParams{
		ID:          bootcampID,
		Name:        pgtype.Text{String: req.Name, Valid: req.Name != ""},
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		StartDate:   startDate,
		EndDate:     endDate,
		IsActive:    pgtype.Bool{Bool: req.IsActive != nil && *req.IsActive, Valid: req.IsActive != nil},
	})
	if err != nil {
		return nil, err
	}

	return s.mapBootcampToData(bootcamp), nil
}

func (s *Service) DeactivateBootcamp(ctx context.Context, bootcampID pgtype.UUID) error {
	return s.queries.ArchiveBootcamp(ctx, bootcampID)
}

// Enrollment operations

func (s *Service) EnrollMember(ctx context.Context, orgID pgtype.UUID, bootcampID pgtype.UUID, req EnrollMemberRequest) (*EnrollmentData, error) {
	memberUUID, err := utils.StringToUUID(req.OrganizationMemberID)
	if err != nil {
		return nil, errors.New("INVALID_MEMBER_ID")
	}

	role, err := s.parseBootcampEnrollmentRole(req.Role)
	if err != nil {
		return nil, err
	}

	// Check if bootcamp exists and is active
	bootcamp, err := s.queries.GetBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, errors.New("BOOTCAMP_NOT_FOUND")
	}

	// Validate bootcamp belongs to the organization
	if bootcamp.OrganizationID != orgID {
		return nil, errors.New("BOOTCAMP_NOT_FOUND")
	}

	if !bootcamp.IsActive {
		return nil, errors.New("BOOTCAMP_INACTIVE")
	}

	// Validate member belongs to the same organization
	orgMember, err := s.queries.GetOrganizationMemberById(ctx, memberUUID)
	if err != nil {
		return nil, errors.New("MEMBER_NOT_FOUND")
	}

	// Check if member belongs to the same organization as the bootcamp
	if orgMember.OrganizationID != bootcamp.OrganizationID {
		return nil, errors.New("CROSS_ORG_VIOLATION")
	}

	enrollment, err := s.queries.EnrollInBootcamp(ctx, db.EnrollInBootcampParams{
		BootcampID:           bootcampID,
		OrganizationMemberID: memberUUID,
		Role:                 role,
		Status:               db.EnrollmentStatusActive,
	})
	if err != nil {
		return nil, err
	}

	return s.mapEnrollmentToData(enrollment), nil
}

func (s *Service) ListEnrollments(ctx context.Context, bootcampID pgtype.UUID) ([]EnrollmentData, error) {
	enrollments, err := s.queries.ListBootcampEnrollments(ctx, bootcampID)
	if err != nil {
		return nil, err
	}

	result := make([]EnrollmentData, len(enrollments))
	for i, enrollment := range enrollments {
		result[i] = EnrollmentData{
			ID:                   enrollment.ID,
			BootcampID:           enrollment.BootcampID,
			OrganizationMemberID: enrollment.OrganizationMemberID,
			Role:                 string(enrollment.Role),
			Status:               string(enrollment.Status),
			EnrolledAt:           enrollment.EnrolledAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			Name:                 enrollment.Name,
			Email:                enrollment.Email.String,
			AvatarUrl:            enrollment.AvatarUrl.String,
			OrgRole:              string(enrollment.OrgRole),
		}
	}

	return result, nil
}

func (s *Service) UpdateEnrollmentRole(ctx context.Context, enrollmentID pgtype.UUID, req UpdateEnrollmentRoleRequest) (*EnrollmentData, error) {
	role, err := s.parseBootcampEnrollmentRole(req.Role)
	if err != nil {
		return nil, err
	}

	enrollment, err := s.queries.UpdateEnrollmentRole(ctx, db.UpdateEnrollmentRoleParams{
		ID:   enrollmentID,
		Role: role,
	})
	if err != nil {
		return nil, err
	}

	return s.mapEnrollmentToData(enrollment), nil
}

func (s *Service) RemoveEnrollment(ctx context.Context, enrollmentID pgtype.UUID) error {
	return s.queries.RemoveEnrollment(ctx, enrollmentID)
}

func (s *Service) GetEnrollment(ctx context.Context, enrollmentID pgtype.UUID) (*EnrollmentData, error) {
	enrollment, err := s.queries.GetEnrollment(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}

	return s.mapEnrollmentToData(enrollment), nil
}

func (s *Service) GetEnrollmentByMember(ctx context.Context, bootcampID pgtype.UUID, memberID pgtype.UUID) (*EnrollmentData, error) {
	enrollment, err := s.queries.GetEnrollmentByMember(ctx, db.GetEnrollmentByMemberParams{
		BootcampID:           bootcampID,
		OrganizationMemberID: memberID,
	})
	if err != nil {
		return nil, err
	}

	return s.mapEnrollmentToData(enrollment), nil
}

// Helper methods

func (s *Service) GetMemberID(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID) (pgtype.UUID, error) {
	member, err := s.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		return pgtype.UUID{}, err
	}
	return member.ID, nil
}

func (s *Service) GetMember(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID) (*db.OrganizationMember, error) {
	member, err := s.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Service) mapBootcampToData(bootcamp db.Bootcamp) *BootcampData {
	return &BootcampData{
		ID:             bootcamp.ID,
		OrganizationID: bootcamp.OrganizationID,
		CreatedBy:      bootcamp.CreatedBy,
		Name:           bootcamp.Name,
		Description:    bootcamp.Description.String,
		StartDate:      FormatDate(bootcamp.StartDate.Time),
		EndDate:        FormatDate(bootcamp.EndDate.Time),
		IsActive:       bootcamp.IsActive,
		CreatedAt:      bootcamp.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      bootcamp.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) mapEnrollmentToData(enrollment db.BootcampEnrollment) *EnrollmentData {
	return &EnrollmentData{
		ID:                   enrollment.ID,
		BootcampID:           enrollment.BootcampID,
		OrganizationMemberID: enrollment.OrganizationMemberID,
		Role:                 string(enrollment.Role),
		Status:               string(enrollment.Status),
		EnrolledAt:           enrollment.EnrolledAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) parseBootcampEnrollmentRole(role string) (db.BootcampEnrollmentRole, error) {
	switch role {
	case "mentor":
		return db.BootcampEnrollmentRoleMentor, nil
	case "mentee":
		return db.BootcampEnrollmentRoleMentee, nil
	default:
		return "", errors.New("INVALID_ROLE")
	}
}
