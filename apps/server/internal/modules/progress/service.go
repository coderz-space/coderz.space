package progress

import (
	"context"
	"errors"

	"github.com/coderz-space/coderz.space/internal/common/utils"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		queries: db.New(pool),
		pool:    pool,
	}
}

// CreateDoubt creates a new doubt for an assignment problem
func (s *Service) CreateDoubt(ctx context.Context, req CreateDoubtRequest, raisedByMemberID pgtype.UUID) (*DoubtData, error) {
	// Parse assignment_problem_id
	assignmentProblemID, err := utils.StringToUUID(req.AssignmentProblemID)
	if err != nil {
		return nil, errors.New("INVALID_ASSIGNMENT_PROBLEM_ID")
	}

	// Validate assignment_problem_id exists and get details
	apDetails, err := s.queries.GetAssignmentProblemDetails(ctx, assignmentProblemID)
	if err != nil {
		return nil, errors.New("ASSIGNMENT_PROBLEM_NOT_FOUND")
	}

	// Verify the problem is assigned to the requesting mentee
	if apDetails.OrganizationMemberID != raisedByMemberID {
		return nil, errors.New("ASSIGNMENT_PROBLEM_NOT_OWNED")
	}

	// Create doubt with raised_by from enrollment context
	doubt, err := s.queries.CreateDoubt(ctx, db.CreateDoubtParams{
		AssignmentProblemID: assignmentProblemID,
		RaisedBy:            raisedByMemberID,
		Message:             req.Message,
	})
	if err != nil {
		return nil, err
	}

	// Get doubt with details for response
	doubtWithDetails, err := s.queries.GetDoubtWithDetails(ctx, doubt.ID)
	if err != nil {
		return nil, err
	}

	return mapDoubtWithDetailsToData(&doubtWithDetails), nil
}

// ListDoubts retrieves doubts with filtering and cursor-based pagination
func (s *Service) ListDoubts(ctx context.Context, bootcampID pgtype.UUID, filters map[string]string, limit int, cursor, userRole string, memberID pgtype.UUID) ([]DoubtData, *CursorPagination, error) {
	// Parse filters
	var assignmentProblemID pgtype.UUID
	if apID, ok := filters["assignment_problem_id"]; ok && apID != "" {
		parsed, err := utils.StringToUUID(apID)
		if err == nil {
			assignmentProblemID = parsed
		}
	}

	var resolved *bool
	if resolvedStr, ok := filters["resolved"]; ok && resolvedStr != "" {
		val := resolvedStr == "true"
		resolved = &val
	}

	// Parse cursor
	var cursorID pgtype.UUID
	if cursor != "" {
		cursorData, err := DecodeCursor(cursor)
		if err == nil && cursorData != nil {
			parsed, err := utils.StringToUUID(cursorData.ID)
			if err == nil {
				cursorID = parsed
			}
		}
	}

	// Fetch doubts with limit + 1 for pagination
	fetchLimit := limit + 1

	var doubts []db.ListDoubtsByMenteeCursorRow
	var err error

	// Role-based filtering
	if userRole == "mentee" {
		// Mentees see only their own doubts
		doubts, err = s.queries.ListDoubtsByMenteeCursor(ctx, db.ListDoubtsByMenteeCursorParams{
			RaisedBy: memberID,
			Column2:  resolved != nil && *resolved,
			Column3:  cursorID,
			Limit:    int32(fetchLimit), // #nosec G115 - fetchLimit is bounded by limit which is max 100
		})
		if err != nil {
			return nil, nil, err
		}

		// Map to common format
		data := make([]DoubtData, 0, len(doubts))
		for i := range doubts {
			data = append(data, mapDoubtMenteeCursorRowToData(&doubts[i]))
		}

		// Check if there are more results
		hasMore := len(doubts) > limit
		if hasMore {
			data = data[:limit]
			doubts = doubts[:limit]
		}

		// Generate next cursor if more results exist
		var nextCursor string
		if hasMore && len(doubts) > 0 {
			lastDoubt := doubts[len(doubts)-1]
			cursor, err := EncodeCursor(lastDoubt.ID, lastDoubt.CreatedAt.Time)
			if err == nil {
				nextCursor = cursor
			}
		}

		pagination := &CursorPagination{
			NextCursor: nextCursor,
			HasMore:    hasMore,
			Limit:      limit,
		}

		return data, pagination, nil
	}

	// Mentors/admins see organization-level doubts
	doubtsCursor, err := s.queries.ListDoubtsCursor(ctx, db.ListDoubtsCursorParams{
		BootcampID: bootcampID,
		Column2:    assignmentProblemID,
		Column3:    resolved != nil && *resolved,
		Column4:    cursorID,
		Limit:      int32(fetchLimit), // #nosec G115 - fetchLimit is bounded by limit which is max 100
	})

	if err != nil {
		return nil, nil, err
	}

	// Check if there are more results
	hasMore := len(doubtsCursor) > limit
	if hasMore {
		doubtsCursor = doubtsCursor[:limit]
	}

	// Generate next cursor if more results exist
	var nextCursor string
	if hasMore && len(doubtsCursor) > 0 {
		lastDoubt := doubtsCursor[len(doubtsCursor)-1]
		cursor, err := EncodeCursor(lastDoubt.ID, lastDoubt.CreatedAt.Time)
		if err == nil {
			nextCursor = cursor
		}
	}

	// Map to response data
	data := make([]DoubtData, len(doubtsCursor))
	for i := range doubtsCursor {
		data[i] = mapDoubtCursorRowToData(&doubtsCursor[i])
	}

	pagination := &CursorPagination{
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	}

	return data, pagination, nil
}

// GetDoubt retrieves a single doubt by ID
func (s *Service) GetDoubt(ctx context.Context, doubtID pgtype.UUID, userRole string, memberID pgtype.UUID) (*DoubtData, error) {
	// Fetch doubt with details
	doubt, err := s.queries.GetDoubtWithDetails(ctx, doubtID)
	if err != nil {
		return nil, errors.New("DOUBT_NOT_FOUND")
	}

	// Validate access based on role (mentees only own, mentors all)
	if userRole == "mentee" && doubt.RaisedBy != memberID {
		return nil, errors.New("ACCESS_DENIED")
	}

	return mapDoubtWithDetailsToData(&doubt), nil
}

// ResolveDoubt marks a doubt as resolved
func (s *Service) ResolveDoubt(ctx context.Context, doubtID, resolvedByMemberID pgtype.UUID, resolutionNote string) (*DoubtData, error) {
	// Fetch doubt to validate it exists
	existingDoubt, err := s.queries.GetDoubt(ctx, doubtID)
	if err != nil {
		return nil, errors.New("DOUBT_NOT_FOUND")
	}

	// Allow idempotent resolution (already resolved is OK)
	if existingDoubt.Resolved {
		// Return existing resolved doubt
		doubtWithDetails, err := s.queries.GetDoubtWithDetails(ctx, doubtID)
		if err != nil {
			return nil, err
		}
		return mapDoubtWithDetailsToData(&doubtWithDetails), nil
	}

	// Validate resolver belongs to same organization
	isSameOrg, err := s.queries.ValidateDoubtResolverOrg(ctx, db.ValidateDoubtResolverOrgParams{
		ResolverMemberID: resolvedByMemberID,
		DoubtID:          doubtID,
	})
	if err != nil {
		return nil, err
	}
	if !isSameOrg {
		return nil, errors.New("RESOLVER_NOT_IN_SAME_ORG")
	}

	// Set resolved to true, resolved_by, resolved_at
	doubt, err := s.queries.ResolveDoubt(ctx, db.ResolveDoubtParams{
		ID:             doubtID,
		ResolvedBy:     resolvedByMemberID,
		ResolutionNote: pgtype.Text{String: resolutionNote, Valid: resolutionNote != ""},
	})
	if err != nil {
		return nil, err
	}

	// Get doubt with details for response
	doubtWithDetails, err := s.queries.GetDoubtWithDetails(ctx, doubt.ID)
	if err != nil {
		return nil, err
	}

	return mapDoubtWithDetailsToData(&doubtWithDetails), nil
}

// DeleteDoubt removes a doubt permanently
func (s *Service) DeleteDoubt(ctx context.Context, doubtID pgtype.UUID, userRole string) error {
	// Validate doubt exists
	_, err := s.queries.GetDoubt(ctx, doubtID)
	if err != nil {
		return errors.New("DOUBT_NOT_FOUND")
	}

	// Enforce only mentors/admins can delete
	if userRole == "mentee" {
		return errors.New("MENTEES_CANNOT_DELETE_DOUBTS")
	}

	// Delete doubt permanently
	return s.queries.DeleteDoubt(ctx, doubtID)
}

// GetMemberIDByUserAndBootcamp retrieves the organization member ID for a user in a bootcamp
func (s *Service) GetMemberIDByUserAndBootcamp(ctx context.Context, userID, bootcampID pgtype.UUID) (pgtype.UUID, error) {
	memberID, err := s.queries.GetMemberIDByUserID(ctx, db.GetMemberIDByUserIDParams{
		UserID:     userID,
		BootcampID: bootcampID,
	})
	if err != nil {
		return pgtype.UUID{}, errors.New("MEMBER_NOT_FOUND")
	}
	return memberID, nil
}

// ValidateAssignmentProblemOwnership verifies that an assignment problem belongs to a mentee
func (s *Service) ValidateAssignmentProblemOwnership(ctx context.Context, assignmentProblemID, memberID pgtype.UUID) error {
	result, err := s.queries.ValidateAssignmentProblemOwnership(ctx, db.ValidateAssignmentProblemOwnershipParams{
		ID:                   assignmentProblemID,
		OrganizationMemberID: memberID,
	})
	if err != nil {
		return err
	}

	if !result {
		return errors.New("ASSIGNMENT_PROBLEM_NOT_OWNED")
	}

	return nil
}

// Helper mapping functions

func mapDoubtWithDetailsToData(d *db.GetDoubtWithDetailsRow) *DoubtData {
	return &DoubtData{
		ID:                  d.ID,
		AssignmentProblemID: d.AssignmentProblemID,
		RaisedBy:            d.RaisedBy,
		Message:             d.Message,
		Resolved:            d.Resolved,
		ResolvedBy:          d.ResolvedBy,
		ResolvedAt:          utils.FormatOptionalTimestamp(d.ResolvedAt),
		ResolutionNote:      formatNullableText(d.ResolutionNote),
		CreatedAt:           utils.FormatTimestamp(d.CreatedAt),
		UpdatedAt:           utils.FormatTimestamp(d.UpdatedAt),
		RaisedByName:        d.RaisedByName,
		RaisedByEmail:       formatNullableText(d.RaisedByEmail),
		ResolvedByName:      formatNullableText(d.ResolvedByName),
	}
}

func mapDoubtCursorRowToData(d *db.ListDoubtsCursorRow) DoubtData {
	return DoubtData{
		ID:                  d.ID,
		AssignmentProblemID: d.AssignmentProblemID,
		RaisedBy:            d.RaisedBy,
		Message:             d.Message,
		Resolved:            d.Resolved,
		ResolvedBy:          d.ResolvedBy,
		ResolvedAt:          utils.FormatOptionalTimestamp(d.ResolvedAt),
		ResolutionNote:      formatNullableText(d.ResolutionNote),
		CreatedAt:           utils.FormatTimestamp(d.CreatedAt),
		UpdatedAt:           utils.FormatTimestamp(d.UpdatedAt),
		RaisedByName:        d.RaisedByName,
		RaisedByEmail:       formatNullableText(d.RaisedByEmail),
		ResolvedByName:      formatNullableText(d.ResolvedByName),
	}
}

func mapDoubtMenteeCursorRowToData(d *db.ListDoubtsByMenteeCursorRow) DoubtData {
	return DoubtData{
		ID:                  d.ID,
		AssignmentProblemID: d.AssignmentProblemID,
		RaisedBy:            d.RaisedBy,
		Message:             d.Message,
		Resolved:            d.Resolved,
		ResolvedBy:          d.ResolvedBy,
		ResolvedAt:          utils.FormatOptionalTimestamp(d.ResolvedAt),
		ResolutionNote:      formatNullableText(d.ResolutionNote),
		CreatedAt:           utils.FormatTimestamp(d.CreatedAt),
		UpdatedAt:           utils.FormatTimestamp(d.UpdatedAt),
		RaisedByName:        d.RaisedByName,
		RaisedByEmail:       formatNullableText(d.RaisedByEmail),
		ResolvedByName:      formatNullableText(d.ResolvedByName),
	}
}

func formatNullableText(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}
