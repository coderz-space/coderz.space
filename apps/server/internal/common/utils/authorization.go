package utils

import (
	"context"
	"errors"

	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuthorizationHelper provides cross-module authorization utilities
type AuthorizationHelper struct {
	queries *db.Queries
}

// NewAuthorizationHelper creates a new authorization helper
func NewAuthorizationHelper(queries *db.Queries) *AuthorizationHelper {
	return &AuthorizationHelper{
		queries: queries,
	}
}

// GetUserOrgMembership retrieves the organization member record for a user in an organization
func (h *AuthorizationHelper) GetUserOrgMembership(ctx context.Context, userID, organizationID pgtype.UUID) (*db.OrganizationMember, error) {
	member, err := h.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, errors.New("USER_NOT_MEMBER_OF_ORGANIZATION")
	}
	return &member, nil
}

// ValidateBootcampAccess verifies that a user has access to a bootcamp
// Returns the bootcamp and user's member ID if access is granted
func (h *AuthorizationHelper) ValidateBootcampAccess(ctx context.Context, userID, bootcampID pgtype.UUID, role string) (*db.Bootcamp, pgtype.UUID, error) {
	// Get bootcamp
	bootcamp, err := h.queries.GetBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, pgtype.UUID{}, errors.New("BOOTCAMP_NOT_FOUND")
	}

	// Get user's organization membership
	member, err := h.GetUserOrgMembership(ctx, userID, bootcamp.OrganizationID)
	if err != nil {
		return nil, pgtype.UUID{}, errors.New("NOT_MEMBER_OF_ORGANIZATION")
	}

	// For mentees, verify they are enrolled in the bootcamp
	if role == "mentee" {
		_, err := h.queries.GetEnrollmentByMember(ctx, db.GetEnrollmentByMemberParams{
			BootcampID:           bootcampID,
			OrganizationMemberID: member.ID,
		})
		if err != nil {
			return nil, pgtype.UUID{}, errors.New("NOT_ENROLLED_IN_BOOTCAMP")
		}
	}

	return &bootcamp, member.ID, nil
}

// ValidateEnrollmentAccess verifies that a user has access to a specific enrollment
// Returns the enrollment if access is granted
func (h *AuthorizationHelper) ValidateEnrollmentAccess(ctx context.Context, userID, enrollmentID pgtype.UUID, role string) (*db.BootcampEnrollment, error) {
	// Get enrollment
	enrollment, err := h.queries.GetEnrollment(ctx, enrollmentID)
	if err != nil {
		return nil, errors.New("ENROLLMENT_NOT_FOUND")
	}

	// Get bootcamp to find organization
	bootcamp, err := h.queries.GetBootcamp(ctx, enrollment.BootcampID)
	if err != nil {
		return nil, errors.New("BOOTCAMP_NOT_FOUND")
	}

	// Get user's organization membership
	member, err := h.GetUserOrgMembership(ctx, userID, bootcamp.OrganizationID)
	if err != nil {
		return nil, errors.New("NOT_MEMBER_OF_ORGANIZATION")
	}

	// For mentees, verify they own this enrollment
	if role == "mentee" && enrollment.OrganizationMemberID != member.ID {
		return nil, errors.New("ACCESS_DENIED")
	}

	return &enrollment, nil
}

// CheckSuperAdmin verifies if a user has super_admin role
func (h *AuthorizationHelper) CheckSuperAdmin(role string) error {
	if role != "super_admin" {
		return errors.New("SUPER_ADMIN_REQUIRED")
	}
	return nil
}

// CheckAdminOrMentor verifies if a user has admin or mentor role
func (h *AuthorizationHelper) CheckAdminOrMentor(role string) error {
	if role != "admin" && role != "mentor" {
		return errors.New("ADMIN_OR_MENTOR_REQUIRED")
	}
	return nil
}

// CheckAdmin verifies if a user has admin role
func (h *AuthorizationHelper) CheckAdmin(role string) error {
	if role != "admin" {
		return errors.New("ADMIN_REQUIRED")
	}
	return nil
}

// ValidateOrgBoundary ensures that a resource belongs to the expected organization
func (h *AuthorizationHelper) ValidateOrgBoundary(_ /* ctx */ context.Context, resourceOrgID, expectedOrgID pgtype.UUID) error {
	if resourceOrgID != expectedOrgID {
		return errors.New("CROSS_ORG_VIOLATION")
	}
	return nil
}

// ValidateBootcampBoundary ensures that a resource belongs to the expected bootcamp
func (h *AuthorizationHelper) ValidateBootcampBoundary(_ /* ctx */ context.Context, resourceBootcampID, expectedBootcampID pgtype.UUID) error {
	if resourceBootcampID != expectedBootcampID {
		return errors.New("CROSS_BOOTCAMP_VIOLATION")
	}
	return nil
}

// ValidateProblemBoundary ensures that a problem belongs to the expected organization
func (h *AuthorizationHelper) ValidateProblemBoundary(ctx context.Context, problemID, expectedOrgID pgtype.UUID) error {
	problem, err := h.queries.GetProblem(ctx, problemID)
	if err != nil {
		return errors.New("PROBLEM_NOT_FOUND")
	}

	if problem.OrganizationID != expectedOrgID {
		return errors.New("CROSS_ORG_VIOLATION")
	}

	return nil
}

// GetMemberIDByUserAndOrg retrieves the organization member ID for a user
func (h *AuthorizationHelper) GetMemberIDByUserAndOrg(ctx context.Context, userID, organizationID pgtype.UUID) (pgtype.UUID, error) {
	member, err := h.GetUserOrgMembership(ctx, userID, organizationID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return member.ID, nil
}

// GetEnrollmentIDByUserAndBootcamp retrieves the enrollment ID for a user in a bootcamp
func (h *AuthorizationHelper) GetEnrollmentIDByUserAndBootcamp(ctx context.Context, userID, bootcampID pgtype.UUID) (pgtype.UUID, error) {
	// Get bootcamp to find organization
	bootcamp, err := h.queries.GetBootcamp(ctx, bootcampID)
	if err != nil {
		return pgtype.UUID{}, errors.New("BOOTCAMP_NOT_FOUND")
	}

	// Get user's organization membership
	member, err := h.GetUserOrgMembership(ctx, userID, bootcamp.OrganizationID)
	if err != nil {
		return pgtype.UUID{}, err
	}

	// Get enrollment
	enrollment, err := h.queries.GetEnrollmentByMember(ctx, db.GetEnrollmentByMemberParams{
		BootcampID:           bootcampID,
		OrganizationMemberID: member.ID,
	})
	if err != nil {
		return pgtype.UUID{}, errors.New("NOT_ENROLLED_IN_BOOTCAMP")
	}

	return enrollment.ID, nil
}
