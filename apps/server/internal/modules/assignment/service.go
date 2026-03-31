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
		Data:    mapAssignmentGroupToData(group),
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

	data := mapAssignmentGroupToData(group)
	data.Problems = make([]GroupProblemRef, len(problems))
	for i, p := range problems {
		data.Problems[i] = GroupProblemRef{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: string(p.Difficulty),
			Position:   p.Position.Int32,
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

	data := mapAssignmentGroupToData(updatedGroup)
	data.Problems = make([]GroupProblemRef, len(problems))
	for i, p := range problems {
		data.Problems[i] = GroupProblemRef{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: string(p.Difficulty),
			Position:   p.Position.Int32,
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
	for i, g := range groups {
		data[i] = mapAssignmentGroupToData(g)
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
	defer tx.Rollback(ctx)

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

	// Calculate deadline if not provided
	var deadlineAt pgtype.Timestamptz
	if req.DeadlineAt != "" {
		t, err := time.Parse(time.RFC3339, req.DeadlineAt)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format: %w", err)
		}
		deadlineAt = pgtype.Timestamptz{Time: t, Valid: true}
	} else {
		// Get group to calculate deadline from deadline_days
		group, err := s.queries.GetAssignmentGroup(ctx, groupID)
		if err != nil {
			return nil, err
		}
		if group.DeadlineDays.Valid {
			deadline := time.Now().Add(time.Duration(group.DeadlineDays.Int32) * 24 * time.Hour)
			deadlineAt = pgtype.Timestamptz{Time: deadline, Valid: true}
		}
	}

	// Use transaction to create assignment and initialize problems
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

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

	// Snapshot problems from group to assignment
	problems, err := qtx.ListAssignmentGroupProblems(ctx, groupID)
	if err != nil {
		return nil, err
	}

	for _, p := range problems {
		_, err := qtx.InitializeAssignmentProblem(ctx, db.InitializeAssignmentProblemParams{
			AssignmentID: assignment.ID,
			ProblemID:    p.ID,
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
		Data:    mapAssignmentToData(assignment),
	}, nil
}

func (s *Service) GetAssignment(ctx context.Context, assignmentID pgtype.UUID) (*AssignmentResponse, error) {
	assignment, err := s.queries.GetAssignment(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	// Get problems for this assignment
	problems, err := s.queries.ListAssignmentProblemsStatus(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	data := mapAssignmentToData(assignment)
	data.Problems = make([]AssignmentProblemData, len(problems))
	for i, p := range problems {
		data.Problems[i] = mapAssignmentProblemToData(p)
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
	for i, a := range assignments {
		assignmentData := AssignmentData{
			ID:                   a.ID,
			AssignmentGroupID:    a.AssignmentGroupID,
			BootcampEnrollmentID: a.BootcampEnrollmentID,
			AssignedBy:           a.AssignedBy,
			AssignedAt:           utils.FormatTimestamp(a.AssignedAt),
			DeadlineAt:           utils.FormatOptionalTimestamp(a.DeadlineAt),
			Status:               string(a.Status),
			CreatedAt:            utils.FormatTimestamp(a.CreatedAt),
			UpdatedAt:            utils.FormatTimestamp(a.UpdatedAt),
			GroupTitle:           a.GroupTitle,
		}
		data[i] = assignmentData
	}

	return &AssignmentListResponse{
		Success: true,
		Data:    data,
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
			Data:    mapAssignmentToData(assignment),
		}, nil
	}

	return nil, fmt.Errorf("no fields to update")
}

// Assignment Problem Progress Methods

func (s *Service) UpdateAssignmentProblemProgress(ctx context.Context, assignmentID, problemID pgtype.UUID, req UpdateAssignmentProblemRequest) (*AssignmentProblemResponse, error) {
	params := db.UpdateAssignmentProblemProgressParams{
		AssignmentID: assignmentID,
		ProblemID:    problemID,
	}

	if req.Status != "" {
		params.Status = db.NullAssignmentProblemStatus{
			AssignmentProblemStatus: db.AssignmentProblemStatus(req.Status),
			Valid:                   true,
		}
		if req.Status == "completed" {
			params.CompletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		}
	}

	if req.SolutionLink != "" {
		params.SolutionLink = pgtype.Text{String: req.SolutionLink, Valid: true}
	}

	if req.Notes != "" {
		params.Notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	problem, err := s.queries.UpdateAssignmentProblemProgress(ctx, params)
	if err != nil {
		return nil, err
	}

	return &AssignmentProblemResponse{
		Success: true,
		Data:    mapAssignmentProblemToDataSimple(problem),
	}, nil
}

func (s *Service) ListAssignmentProblems(ctx context.Context, assignmentID pgtype.UUID) (*AssignmentProblemListResponse, error) {
	problems, err := s.queries.ListAssignmentProblemsStatus(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	data := make([]AssignmentProblemData, len(problems))
	for i, p := range problems {
		data[i] = mapAssignmentProblemToData(p)
	}

	return &AssignmentProblemListResponse{
		Success: true,
		Data:    data,
	}, nil
}

// Helper mapping functions

func mapAssignmentGroupToData(g db.AssignmentGroup) AssignmentGroupData {
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

func mapAssignmentToData(a db.Assignment) AssignmentData {
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

func mapAssignmentProblemToData(p db.ListAssignmentProblemsStatusRow) AssignmentProblemData {
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

func mapAssignmentProblemToDataSimple(p db.AssignmentProblem) AssignmentProblemData {
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
