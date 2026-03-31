package assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewService(pool *pgxpool.Pool, queries *db.Queries) *Service {
	return &Service{
		pool:    pool,
		queries: queries,
	}
}

// Assignment Group Methods

func (s *Service) CreateAssignmentGroup(ctx context.Context, req CreateAssignmentGroupRequest, bootcampID, createdBy pgtype.UUID) (*AssignmentGroupResponse, error) {
	// Validate bootcamp exists and is accessible
	bootcamp, err := s.queries.GetBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, err
	}

	// Verify bootcamp is active
	if !bootcamp.IsActive {
		return nil, fmt.Errorf("BOOTCAMP_INACTIVE")
	}

	// Create assignment group
	group, err := s.queries.CreateAssignmentGroup(ctx, db.CreateAssignmentGroupParams{
		BootcampID:   bootcampID,
		CreatedBy:    createdBy,
		Title:        req.Title,
		Description:  pgtype.Text{String: req.Description, Valid: req.Description != ""},
		DeadlineDays: pgtype.Int4{Int32: req.DeadlineDays, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AssignmentGroupResponse{
		Success: true,
		Data:    mapAssignmentGroupToData(&group),
	}, nil
}

func (s *Service) GetAssignmentGroup(ctx context.Context, groupID pgtype.UUID) (*AssignmentGroupResponse, error) {
	group, err := s.queries.GetAssignmentGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Get problems for this group
	problems, err := s.queries.ListAssignmentGroupProblems(ctx, groupID)
	if err != nil {
		return nil, err
	}

	data := mapAssignmentGroupToData(&group)
	data.Problems = make([]GroupProblemRef, len(problems))
	for i := range problems {
		data.Problems[i] = GroupProblemRef{
			ProblemID:  problems[i].ID,
			Title:      problems[i].Title,
			Difficulty: string(problems[i].Difficulty),
			Position:   problems[i].Position.Int32,
		}
	}

	return &AssignmentGroupResponse{
		Success: true,
		Data:    data,
	}, nil
}

func (s *Service) UpdateAssignmentGroup(ctx context.Context, groupID pgtype.UUID, req UpdateAssignmentGroupRequest) (*AssignmentGroupResponse, error) {
	// Validate at least one field is provided
	if req.Title == "" && req.Description == "" && req.DeadlineDays == 0 {
		return nil, fmt.Errorf("NO_FIELDS_PROVIDED")
	}

	// Get existing group to verify it exists
	existingGroup, err := s.queries.GetAssignmentGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Prepare update parameters
	params := db.UpdateAssignmentGroupParams{
		ID: groupID,
	}

	// Only update fields that are provided
	if req.Title != "" {
		params.Title = pgtype.Text{String: req.Title, Valid: true}
	}
	if req.Description != "" {
		params.Description = pgtype.Text{String: req.Description, Valid: true}
	}
	if req.DeadlineDays > 0 {
		params.DeadlineDays = pgtype.Int4{Int32: req.DeadlineDays, Valid: true}
	}

	// Update the assignment group
	updatedGroup, err := s.queries.UpdateAssignmentGroup(ctx, params)
	if err != nil {
		return nil, err
	}

	// Note: bootcamp_id is immutable and cannot be changed (as per requirements 7.7, 7.8)
	// Existing assignment instances are not modified (as per requirement 7.7)
	_ = existingGroup // Used for validation

	// Get problems for the updated group
	problems, err := s.queries.ListAssignmentGroupProblems(ctx, groupID)
	if err != nil {
		return nil, err
	}

	data := mapAssignmentGroupToData(&updatedGroup)
	data.Problems = make([]GroupProblemRef, len(problems))
	for i := range problems {
		data.Problems[i] = GroupProblemRef{
			ProblemID:  problems[i].ID,
			Title:      problems[i].Title,
			Difficulty: string(problems[i].Difficulty),
			Position:   problems[i].Position.Int32,
		}
	}

	return &AssignmentGroupResponse{
		Success: true,
		Data:    data,
	}, nil
}

func (s *Service) ListAssignmentGroups(ctx context.Context, bootcampID pgtype.UUID, createdBy *pgtype.UUID, page, limit int) (*AssignmentGroupListResponse, error) {
	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Count total groups
	countParams := db.CountAssignmentGroupsByBootcampParams{
		BootcampID: bootcampID,
		CreatedBy:  pgtype.UUID{},
	}
	if createdBy != nil {
		countParams.CreatedBy = *createdBy
	}

	total, err := s.queries.CountAssignmentGroupsByBootcamp(ctx, countParams)
	if err != nil {
		return nil, err
	}

	// List groups with pagination
	listParams := db.ListAssignmentGroupsByBootcampParams{
		BootcampID: bootcampID,
		CreatedBy:  pgtype.UUID{},
		Limit:      int32(limit),  // #nosec G115 - limit is bounded to max 100
		Offset:     int32(offset), // #nosec G115 - offset is calculated from bounded values
	}
	if createdBy != nil {
		listParams.CreatedBy = *createdBy
	}

	groups, err := s.queries.ListAssignmentGroupsByBootcamp(ctx, listParams)
	if err != nil {
		return nil, err
	}

	data := make([]AssignmentGroupData, len(groups))
	for i := range groups {
		data[i] = mapAssignmentGroupToData(&groups[i])
	}

	return &AssignmentGroupListResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	}, nil
}

func (s *Service) AddProblemsToGroup(ctx context.Context, groupID pgtype.UUID, req AddProblemsToGroupRequest) error {
	// Use transaction to ensure atomicity
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) //nolint:errcheck // Rollback is safe to call even after commit
	}()

	qtx := s.queries.WithTx(tx)

	for _, p := range req.Problems {
		problemID, err := utils.StringToUUID(p.ProblemID)
		if err != nil {
			return fmt.Errorf("invalid problem ID: %w", err)
		}

		err = qtx.AddProblemToAssignmentGroup(ctx, db.AddProblemToAssignmentGroupParams{
			AssignmentGroupID: groupID,
			ProblemID:         problemID,
			Position:          pgtype.Int4{Int32: p.Position, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Service) RemoveProblemFromGroup(ctx context.Context, groupID, problemID pgtype.UUID) error {
	return s.queries.RemoveProblemFromAssignmentGroup(ctx, db.RemoveProblemFromAssignmentGroupParams{
		AssignmentGroupID: groupID,
		ProblemID:         problemID,
	})
}

func (s *Service) ReplaceGroupProblems(ctx context.Context, groupID pgtype.UUID, req ReplaceGroupProblemsRequest) error {
	// Validate all problem_ids are unique (Requirement 7.11)
	problemIDSet := make(map[string]bool)
	for _, p := range req.Problems {
		if problemIDSet[p.ProblemID] {
			return fmt.Errorf("DUPLICATE_PROBLEM_ID: problem ID %s appears multiple times", p.ProblemID)
		}
		problemIDSet[p.ProblemID] = true
	}

	// Validate all positions are unique positive integers (Requirement 7.12)
	positionSet := make(map[int32]bool)
	for _, p := range req.Problems {
		if p.Position < 1 {
			return fmt.Errorf("INVALID_POSITION: position must be a positive integer, got %d", p.Position)
		}
		if positionSet[p.Position] {
			return fmt.Errorf("DUPLICATE_POSITION: position %d appears multiple times", p.Position)
		}
		positionSet[p.Position] = true
	}

	// Execute replacement atomically in transaction (Requirement 7.13, 20.9)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback is safe to call even after commit
	}()

	qtx := s.queries.WithTx(tx)

	// Clear all existing problems from the group
	err = qtx.ClearAssignmentGroupProblems(ctx, groupID)
	if err != nil {
		return err
	}

	// Add all new problems
	for _, p := range req.Problems {
		problemID, err := utils.StringToUUID(p.ProblemID)
		if err != nil {
			return fmt.Errorf("invalid problem ID: %w", err)
		}

		err = qtx.AddProblemToAssignmentGroup(ctx, db.AddProblemToAssignmentGroupParams{
			AssignmentGroupID: groupID,
			ProblemID:         problemID,
			Position:          pgtype.Int4{Int32: p.Position, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Service) DeleteAssignmentGroup(ctx context.Context, groupID pgtype.UUID) error {
	// Check if there are any existing assignments for this group
	count, err := s.queries.CountAssignmentsByGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// Return conflict error if assignments exist (Requirements 7.9, 25.7)
	if count > 0 {
		return fmt.Errorf("ASSIGNMENT_GROUP_HAS_ASSIGNMENTS")
	}

	// Delete the assignment group
	return s.queries.DeleteAssignmentGroup(ctx, groupID)
}

// Assignment Instance Methods

func (s *Service) CreateAssignment(ctx context.Context, req CreateAssignmentRequest, assignedBy pgtype.UUID) (*AssignmentResponse, error) {
	groupID, err := utils.StringToUUID(req.AssignmentGroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid assignment group ID: %w", err)
	}

	enrollmentID, err := utils.StringToUUID(req.BootcampEnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid bootcamp enrollment ID: %w", err)
	}

	// Validate assignment_group_id exists (Requirement 8.1)
	group, err := s.queries.GetAssignmentGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("ASSIGNMENT_GROUP_NOT_FOUND")
	}

	// Validate bootcamp_enrollment_id belongs to same bootcamp (Requirement 8.2)
	enrollmentBootcamp, err := s.queries.GetEnrollmentBootcamp(ctx, enrollmentID)
	if err != nil {
		return nil, fmt.Errorf("ENROLLMENT_NOT_FOUND")
	}

	if enrollmentBootcamp.BootcampID != group.BootcampID {
		return nil, fmt.Errorf("ENROLLMENT_BOOTCAMP_MISMATCH")
	}

	// Validate enrollment is active (Requirement 8.5)
	if !enrollmentBootcamp.IsActive {
		return nil, fmt.Errorf("BOOTCAMP_INACTIVE")
	}

	// Prevent duplicate active assignments (Requirement 8.4)
	duplicateCount, err := s.queries.CheckDuplicateActiveAssignment(ctx, db.CheckDuplicateActiveAssignmentParams{
		AssignmentGroupID:    groupID,
		BootcampEnrollmentID: enrollmentID,
	})
	if err != nil {
		return nil, err
	}
	if duplicateCount > 0 {
		return nil, fmt.Errorf("DUPLICATE_ACTIVE_ASSIGNMENT")
	}

	// Calculate deadline_at from assigned_at + deadline_days if not provided (Requirement 8.3)
	var deadlineAt pgtype.Timestamptz
	if req.DeadlineAt != "" {
		t, err := time.Parse(time.RFC3339, req.DeadlineAt)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format: %w", err)
		}
		deadlineAt = pgtype.Timestamptz{Time: t, Valid: true}
	} else if group.DeadlineDays.Valid {
		deadline := time.Now().Add(time.Duration(group.DeadlineDays.Int32) * 24 * time.Hour)
		deadlineAt = pgtype.Timestamptz{Time: deadline, Valid: true}
	}

	// Use transaction to create assignment and snapshot problems atomically (Requirements 8.7, 28.1, 28.2, 28.3, 28.7)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback is safe to call even after commit
	}()

	qtx := s.queries.WithTx(tx)

	// Create assignment instance
	assignment, err := qtx.AssignGroupToMentee(ctx, db.AssignGroupToMenteeParams{
		AssignmentGroupID:    groupID,
		BootcampEnrollmentID: enrollmentID,
		AssignedBy:           assignedBy,
		DeadlineAt:           deadlineAt,
		Status:               "active",
	})
	if err != nil {
		return nil, err
	}

	// Snapshot group problems into assignment_problems atomically (Requirements 28.1, 28.6)
	problems, err := qtx.ListAssignmentGroupProblems(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Initialize all problems with pending status (Requirement 28.6)
	for i := range problems {
		_, err := qtx.InitializeAssignmentProblem(ctx, db.InitializeAssignmentProblemParams{
			AssignmentID: assignment.ID,
			ProblemID:    problems[i].ID,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &AssignmentResponse{
		Success: true,
		Data:    mapAssignmentToData(&assignment),
	}, nil
}

func (s *Service) GetAssignment(ctx context.Context, assignmentID pgtype.UUID) (*AssignmentResponse, error) {
	// Get assignment with group metadata (Requirement 8.13)
	assignment, err := s.queries.GetAssignmentWithGroup(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	// Get problems for this assignment
	problems, err := s.queries.ListAssignmentProblemsStatus(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	data := AssignmentData{
		ID:                   assignment.ID,
		AssignmentGroupID:    assignment.AssignmentGroupID,
		BootcampEnrollmentID: assignment.BootcampEnrollmentID,
		AssignedBy:           assignment.AssignedBy,
		AssignedAt:           utils.FormatTimestamp(assignment.AssignedAt),
		DeadlineAt:           utils.FormatOptionalTimestamp(assignment.DeadlineAt),
		Status:               string(assignment.Status),
		CreatedAt:            utils.FormatTimestamp(assignment.CreatedAt),
		UpdatedAt:            utils.FormatTimestamp(assignment.UpdatedAt),
		GroupTitle:           assignment.GroupTitle,
	}

	data.Problems = make([]AssignmentProblemData, len(problems))
	for i := range problems {
		data.Problems[i] = mapAssignmentProblemToData(&problems[i])
	}

	return &AssignmentResponse{
		Success: true,
		Data:    data,
	}, nil
}

func (s *Service) ListAssignmentsByMentee(ctx context.Context, enrollmentID pgtype.UUID) (*AssignmentListResponse, error) {
	assignments, err := s.queries.ListAssignmentsByMentee(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}

	data := make([]AssignmentData, len(assignments))
	for i := range assignments {
		data[i] = AssignmentData{
			ID:                   assignments[i].ID,
			AssignmentGroupID:    assignments[i].AssignmentGroupID,
			BootcampEnrollmentID: assignments[i].BootcampEnrollmentID,
			AssignedBy:           assignments[i].AssignedBy,
			AssignedAt:           utils.FormatTimestamp(assignments[i].AssignedAt),
			DeadlineAt:           utils.FormatOptionalTimestamp(assignments[i].DeadlineAt),
			Status:               string(assignments[i].Status),
			CreatedAt:            utils.FormatTimestamp(assignments[i].CreatedAt),
			UpdatedAt:            utils.FormatTimestamp(assignments[i].UpdatedAt),
			GroupTitle:           assignments[i].GroupTitle,
		}
	}

	return &AssignmentListResponse{
		Success: true,
		Data:    data,
	}, nil
}

// ListAssignments returns assignments with filtering by bootcamp_id, assignment_group_id, and status
// Supports pagination (Requirement 8.8, 8.9, 23.1)
func (s *Service) ListAssignments(ctx context.Context, bootcampID pgtype.UUID, assignmentGroupID *pgtype.UUID, status *string, page, limit int) (*AssignmentListResponse, error) {
	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Count total assignments
	countParams := db.CountAssignmentsParams{
		BootcampID:        bootcampID,
		AssignmentGroupID: pgtype.UUID{},
		Status:            db.NullAssignmentStatus{},
	}
	if assignmentGroupID != nil {
		countParams.AssignmentGroupID = *assignmentGroupID
	}
	if status != nil {
		countParams.Status = db.NullAssignmentStatus{
			AssignmentStatus: db.AssignmentStatus(*status),
			Valid:            true,
		}
	}

	total, err := s.queries.CountAssignments(ctx, countParams)
	if err != nil {
		return nil, err
	}

	// List assignments with pagination
	listParams := db.ListAssignmentsParams{
		BootcampID:        bootcampID,
		AssignmentGroupID: pgtype.UUID{},
		Status:            db.NullAssignmentStatus{},
		Limit:             int32(limit),  // #nosec G115 - limit is bounded to max 100
		Offset:            int32(offset), // #nosec G115 - offset is calculated from bounded values
	}
	if assignmentGroupID != nil {
		listParams.AssignmentGroupID = *assignmentGroupID
	}
	if status != nil {
		listParams.Status = db.NullAssignmentStatus{
			AssignmentStatus: db.AssignmentStatus(*status),
			Valid:            true,
		}
	}

	assignments, err := s.queries.ListAssignments(ctx, listParams)
	if err != nil {
		return nil, err
	}

	data := make([]AssignmentData, len(assignments))
	for i := range assignments {
		data[i] = AssignmentData{
			ID:                   assignments[i].ID,
			AssignmentGroupID:    assignments[i].AssignmentGroupID,
			BootcampEnrollmentID: assignments[i].BootcampEnrollmentID,
			AssignedBy:           assignments[i].AssignedBy,
			AssignedAt:           utils.FormatTimestamp(assignments[i].AssignedAt),
			DeadlineAt:           utils.FormatOptionalTimestamp(assignments[i].DeadlineAt),
			Status:               string(assignments[i].Status),
			CreatedAt:            utils.FormatTimestamp(assignments[i].CreatedAt),
			UpdatedAt:            utils.FormatTimestamp(assignments[i].UpdatedAt),
			GroupTitle:           assignments[i].GroupTitle,
		}
	}

	return &AssignmentListResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	}, nil
}

func (s *Service) UpdateAssignment(ctx context.Context, assignmentID pgtype.UUID, req UpdateAssignmentRequest) (*AssignmentResponse, error) {
	// For now, only support status updates
	if req.Status != "" {
		assignment, err := s.queries.UpdateAssignmentStatus(ctx, db.UpdateAssignmentStatusParams{
			ID:     assignmentID,
			Status: db.AssignmentStatus(req.Status),
		})
		if err != nil {
			return nil, err
		}

		return &AssignmentResponse{
			Success: true,
			Data:    mapAssignmentToData(&assignment),
		}, nil
	}

	return nil, fmt.Errorf("no fields to update")
}

// UpdateAssignmentDeadline updates the deadline of an assignment (Requirement 8.10)
func (s *Service) UpdateAssignmentDeadline(ctx context.Context, assignmentID pgtype.UUID, req UpdateAssignmentDeadlineRequest) (*AssignmentResponse, error) {
	// Validate new deadline is valid timestamp (Requirement 8.10)
	t, err := time.Parse(time.RFC3339, req.DeadlineAt)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DEADLINE_FORMAT")
	}

	assignment, err := s.queries.UpdateAssignmentDeadline(ctx, db.UpdateAssignmentDeadlineParams{
		ID:         assignmentID,
		DeadlineAt: pgtype.Timestamptz{Time: t, Valid: true},
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("ASSIGNMENT_NOT_FOUND")
		}
		return nil, err
	}

	return &AssignmentResponse{
		Success: true,
		Data:    mapAssignmentToData(&assignment),
	}, nil
}

// UpdateAssignmentStatus updates the status of an assignment (Requirement 8.11)
func (s *Service) UpdateAssignmentStatus(ctx context.Context, assignmentID pgtype.UUID, req UpdateAssignmentStatusRequest) (*AssignmentResponse, error) {
	// Validate status transitions (Requirement 8.11)
	validStatuses := map[string]bool{
		"active":    true,
		"completed": true,
		"expired":   true,
	}

	if !validStatuses[req.Status] {
		return nil, fmt.Errorf("INVALID_STATUS")
	}

	assignment, err := s.queries.UpdateAssignmentStatus(ctx, db.UpdateAssignmentStatusParams{
		ID:     assignmentID,
		Status: db.AssignmentStatus(req.Status),
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("ASSIGNMENT_NOT_FOUND")
		}
		return nil, err
	}

	return &AssignmentResponse{
		Success: true,
		Data:    mapAssignmentToData(&assignment),
	}, nil
}

// Assignment Problem Progress Methods

func (s *Service) UpdateAssignmentProblemProgress(ctx context.Context, assignmentID, problemID pgtype.UUID, req UpdateAssignmentProblemRequest, userID pgtype.UUID) (*AssignmentProblemResponse, error) {
	// Get assignment with enrollment to verify ownership
	_, err := s.queries.GetAssignmentWithEnrollment(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found")
	}

	// TODO: Verify mentee owns the assignment by checking organization_member_id matches user
	// This requires additional query to map user_id to organization_member_id

	// Get current problem status to check for regression
	currentProblem, err := s.queries.GetAssignmentProblem(ctx, db.GetAssignmentProblemParams{
		AssignmentID: assignmentID,
		ProblemID:    problemID,
	})
	if err != nil {
		return nil, fmt.Errorf("problem not found in assignment")
	}

	// Prevent status regression from completed to pending (Requirement 9.10)
	if currentProblem.Status == db.AssignmentProblemStatusCompleted && req.Status == "pending" {
		return nil, fmt.Errorf("cannot regress status from completed to pending")
	}

	params := db.UpdateAssignmentProblemProgressParams{
		AssignmentID: assignmentID,
		ProblemID:    problemID,
	}

	// Validate and set status (Requirement 9.1)
	if req.Status != "" {
		// Status validation is already handled by the DTO validation tag
		params.Status = db.NullAssignmentProblemStatus{
			AssignmentProblemStatus: db.AssignmentProblemStatus(req.Status),
			Valid:                   true,
		}
		// Set completed_at when status changes to completed (Requirement 9.2)
		if req.Status == "completed" {
			params.CompletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		}
	}

	// Validate and set solution_link (Requirement 9.3)
	if req.SolutionLink != "" {
		// URL validation is already handled by the DTO validation tag
		params.SolutionLink = pgtype.Text{String: req.SolutionLink, Valid: true}
	}

	// Set notes if provided
	if req.Notes != "" {
		params.Notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	problem, err := s.queries.UpdateAssignmentProblemProgress(ctx, params)
	if err != nil {
		return nil, err
	}

	return &AssignmentProblemResponse{
		Success: true,
		Data:    mapAssignmentProblemToDataSimple(&problem),
	}, nil
}

func (s *Service) ListAssignmentProblems(ctx context.Context, assignmentID pgtype.UUID) (*AssignmentProblemListResponse, error) {
	problems, err := s.queries.ListAssignmentProblemsStatus(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	data := make([]AssignmentProblemData, len(problems))
	for i := range problems {
		data[i] = mapAssignmentProblemToData(&problems[i])
	}

	return &AssignmentProblemListResponse{
		Success: true,
		Data:    data,
	}, nil
}

// GetAssignmentProblem retrieves a single assignment problem with details (Requirement 9.8, 9.9)
func (s *Service) GetAssignmentProblem(ctx context.Context, assignmentID, problemID pgtype.UUID) (*AssignmentProblemResponse, error) {
	problem, err := s.queries.GetAssignmentProblem(ctx, db.GetAssignmentProblemParams{
		AssignmentID: assignmentID,
		ProblemID:    problemID,
	})
	if err != nil {
		return nil, fmt.Errorf("problem not found in assignment")
	}

	return &AssignmentProblemResponse{
		Success: true,
		Data:    mapGetAssignmentProblemToData(&problem),
	}, nil
}

// Helper mapping functions

func mapAssignmentGroupToData(g *db.AssignmentGroup) AssignmentGroupData {
	return AssignmentGroupData{
		ID:           g.ID,
		BootcampID:   g.BootcampID,
		CreatedBy:    g.CreatedBy,
		Title:        g.Title,
		Description:  g.Description.String,
		DeadlineDays: g.DeadlineDays.Int32,
		CreatedAt:    utils.FormatTimestamp(g.CreatedAt),
		UpdatedAt:    utils.FormatTimestamp(g.UpdatedAt),
	}
}

func mapAssignmentToData(a *db.Assignment) AssignmentData {
	return AssignmentData{
		ID:                   a.ID,
		AssignmentGroupID:    a.AssignmentGroupID,
		BootcampEnrollmentID: a.BootcampEnrollmentID,
		AssignedBy:           a.AssignedBy,
		AssignedAt:           utils.FormatTimestamp(a.AssignedAt),
		DeadlineAt:           utils.FormatOptionalTimestamp(a.DeadlineAt),
		Status:               string(a.Status),
		CreatedAt:            utils.FormatTimestamp(a.CreatedAt),
		UpdatedAt:            utils.FormatTimestamp(a.UpdatedAt),
	}
}

func mapAssignmentProblemToData(p *db.ListAssignmentProblemsStatusRow) AssignmentProblemData {
	return AssignmentProblemData{
		ID:           p.ID,
		AssignmentID: p.AssignmentID,
		ProblemID:    p.ProblemID,
		Status:       string(p.Status),
		SolutionLink: p.SolutionLink.String,
		Notes:        p.Notes.String,
		CompletedAt:  utils.FormatOptionalTimestamp(p.CompletedAt),
		CreatedAt:    utils.FormatTimestamp(p.CreatedAt),
		UpdatedAt:    utils.FormatTimestamp(p.UpdatedAt),
		Title:        p.Title,
		Difficulty:   string(p.Difficulty),
	}
}

func mapGetAssignmentProblemToData(p *db.GetAssignmentProblemRow) AssignmentProblemData {
	return AssignmentProblemData{
		ID:           p.ID,
		AssignmentID: p.AssignmentID,
		ProblemID:    p.ProblemID,
		Status:       string(p.Status),
		SolutionLink: p.SolutionLink.String,
		Notes:        p.Notes.String,
		CompletedAt:  utils.FormatOptionalTimestamp(p.CompletedAt),
		CreatedAt:    utils.FormatTimestamp(p.CreatedAt),
		UpdatedAt:    utils.FormatTimestamp(p.UpdatedAt),
		Title:        p.Title,
		Difficulty:   string(p.Difficulty),
	}
}

func mapAssignmentProblemToDataSimple(p *db.AssignmentProblem) AssignmentProblemData {
	return AssignmentProblemData{
		ID:           p.ID,
		AssignmentID: p.AssignmentID,
		ProblemID:    p.ProblemID,
		Status:       string(p.Status),
		SolutionLink: p.SolutionLink.String,
		Notes:        p.Notes.String,
		CompletedAt:  utils.FormatOptionalTimestamp(p.CompletedAt),
		CreatedAt:    utils.FormatTimestamp(p.CreatedAt),
		UpdatedAt:    utils.FormatTimestamp(p.UpdatedAt),
	}
}
