package mentorship

import (
	"context"
	"errors"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) CreateRoleRequest(ctx context.Context, userIDStr string, role string) error {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return err
	}

	// 1. Get or create a default organization
	orgs, err := s.queries.GetPendingOrganizations(ctx) // fallback logic for MVP
	if err != nil {
		return err
	}
	if len(orgs) == 0 {
		return errors.New("no organization found to join")
	}
	orgID := orgs[0].ID

	// 2. Add user to organization as member
	member, err := s.queries.AddOrganizationMember(ctx, db.AddOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           db.OrgMemberRole(role),
	})
	if err != nil {
		// Attempting to fallback if already exists:
		return errors.New("could not create member request")
	}

	// 3. Get or create default bootcamp
	bootcamps, err := s.queries.ListBootcampsByOrg(ctx, orgID)
	if err != nil || len(bootcamps) == 0 {
		return errors.New("no bootcamp found")
	}

	// 4. Enroll with 'pending' status
	_, err = s.queries.EnrollInBootcamp(ctx, db.EnrollInBootcampParams{
		BootcampID:           bootcamps[0].ID,
		OrganizationMemberID: member.ID,
		Role:                 db.BootcampEnrollmentRole(role),
		Status:               "pending",
	})
	return err
}

func (s *Service) ListPendingRequests(ctx context.Context) ([]MenteeRequestDTO, error) {
	rows, err := s.queries.WebListPendingRequests(ctx)
	if err != nil {
		return nil, err
	}

	var res []MenteeRequestDTO
	for _, r := range rows {
		res = append(res, MenteeRequestDTO{
			ID:         utils.UUIDToString(r.ID),
			FirstName:  r.FirstName,
			Email:      r.Email.String,
			SignedUpAt: r.SignedUpAt.Time.Format("2006-01-02T15:04:05Z"),
			Status:     string(r.Status),
		})
	}
	return res, nil
}

func (s *Service) UpdateStatus(ctx context.Context, idStr string, status string) error {
	id, err := utils.StringToUUID(idStr)
	if err != nil {
		return err
	}
	
	dbStatus := db.EnrollmentStatus(status)
	_, err = s.queries.UpdateEnrollmentStatus(ctx, db.UpdateEnrollmentStatusParams{
		ID:     id,
		Status: dbStatus,
	})
	return err
}

func (s *Service) DeleteRequest(ctx context.Context, idStr string) error {
	id, err := utils.StringToUUID(idStr)
	if err != nil {
		return err
	}
	return s.queries.RemoveEnrollment(ctx, id)
}
