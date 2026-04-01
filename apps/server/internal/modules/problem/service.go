package problem

import (
	"context"
	"errors"

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

// Problem operations

func (s *Service) CreateProblem(ctx context.Context, req CreateProblemRequest, orgID, userID pgtype.UUID) (*ProblemData, error) {
	// Verify user is a member of the organization
	member, err := s.GetMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	// Create problem
	problem, err := s.queries.CreateProblem(ctx, db.CreateProblemParams{
		OrganizationID: orgID,
		CreatedBy:      member.ID,
		Title:          req.Title,
		Description:    pgtype.Text{String: req.Description, Valid: true},
		Difficulty:     db.DifficultyLevel(req.Difficulty),
		ExternalLink:   pgtype.Text{String: req.ExternalLink, Valid: req.ExternalLink != ""},
	})
	if err != nil {
		return nil, err
	}

	return s.mapProblemToData(&problem), nil
}

func (s *Service) ListProblems(ctx context.Context, orgID pgtype.UUID) ([]ProblemData, error) {
	problems, err := s.queries.ListProblemsByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	result := make([]ProblemData, len(problems))
	for i := range problems {
		result[i] = *s.mapProblemToData(&problems[i])
	}

	return result, nil
}

func (s *Service) GetProblem(ctx context.Context, problemID pgtype.UUID) (*ProblemData, error) {
	problem, err := s.queries.GetProblem(ctx, problemID)
	if err != nil {
		return nil, err
	}

	data := s.mapProblemToData(&problem)

	// Load tags
	tags, err := s.queries.ListProblemTags(ctx, problemID)
	if err == nil && len(tags) > 0 {
		data.Tags = make([]TagData, len(tags))
		for i, tag := range tags {
			data.Tags[i] = *s.mapTagToData(&tag)
		}
	}

	// Load resources
	resources, err := s.queries.ListProblemResources(ctx, problemID)
	if err == nil && len(resources) > 0 {
		data.Resources = make([]ResourceData, len(resources))
		for i, resource := range resources {
			data.Resources[i] = *s.mapResourceToData(&resource)
		}
	}

	return data, nil
}

func (s *Service) UpdateProblem(ctx context.Context, req UpdateProblemRequest, problemID pgtype.UUID) (*ProblemData, error) {
	// Check if at least one field is provided
	if req.Title == "" && req.Description == "" && req.Difficulty == "" && req.ExternalLink == "" {
		return nil, errors.New("NO_FIELDS_PROVIDED")
	}

	// Build update params
	params := db.UpdateProblemParams{
		ID: problemID,
	}

	if req.Title != "" {
		params.Title = pgtype.Text{String: req.Title, Valid: true}
	}
	if req.Description != "" {
		params.Description = pgtype.Text{String: req.Description, Valid: true}
	}
	if req.Difficulty != "" {
		params.Difficulty = db.NullDifficultyLevel{
			DifficultyLevel: db.DifficultyLevel(req.Difficulty),
			Valid:           true,
		}
	}
	if req.ExternalLink != "" {
		params.ExternalLink = pgtype.Text{String: req.ExternalLink, Valid: true}
	}

	problem, err := s.queries.UpdateProblem(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.mapProblemToData(&problem), nil
}

func (s *Service) DeleteProblem(ctx context.Context, problemID pgtype.UUID) error {
	count, err := s.queries.CountProblemAssignments(ctx, problemID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("PROBLEM_REFERENCED_BY_ASSIGNMENTS")
	}

	return s.queries.ArchiveProblem(ctx, problemID)
}

// Tag operations

func (s *Service) CreateTag(ctx context.Context, req CreateTagRequest, orgID, userID pgtype.UUID) (*TagData, error) {
	// Verify user is a member of the organization
	member, err := s.GetMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	// Normalize tag name
	normalizedName := NormalizeTagName(req.Name)

	// Check if tag already exists
	_, err = s.queries.GetTagByName(ctx, db.GetTagByNameParams{
		OrganizationID: orgID,
		Name:           normalizedName,
	})
	if err == nil {
		// Tag already exists, return conflict error
		return nil, errors.New("TAG_ALREADY_EXISTS")
	}

	// Create tag
	tag, err := s.queries.CreateTag(ctx, db.CreateTagParams{
		OrganizationID: orgID,
		CreatedBy:      member.ID,
		Name:           normalizedName,
	})
	if err != nil {
		return nil, err
	}

	return s.mapTagToData(&tag), nil
}

func (s *Service) ListTags(ctx context.Context, orgID pgtype.UUID, searchQuery string) ([]TagData, error) {
	var tags []db.Tag
	var err error

	if searchQuery != "" {
		tags, err = s.queries.SearchTagsByName(ctx, db.SearchTagsByNameParams{
			OrganizationID: orgID,
			Name:           searchQuery,
		})
	} else {
		tags, err = s.queries.ListTagsByOrg(ctx, orgID)
	}

	if err != nil {
		return nil, err
	}

	result := make([]TagData, len(tags))
	for i := range tags {
		result[i] = *s.mapTagToData(&tags[i])
	}

	return result, nil
}

func (s *Service) GetTag(ctx context.Context, tagID pgtype.UUID) (*db.Tag, error) {
	tag, err := s.queries.GetTag(ctx, tagID)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *Service) UpdateTag(ctx context.Context, req UpdateTagRequest, tagID, orgID pgtype.UUID) (*TagData, error) {
	// Normalize new tag name
	normalizedName := NormalizeTagName(req.Name)

	// Check if tag with new name already exists (excluding current tag)
	existingTag, err := s.queries.GetTagByName(ctx, db.GetTagByNameParams{
		OrganizationID: orgID,
		Name:           normalizedName,
	})
	if err == nil && existingTag.ID.Bytes != tagID.Bytes {
		// Another tag with this name already exists
		return nil, errors.New("TAG_NAME_ALREADY_EXISTS")
	}

	// Update tag
	tag, err := s.queries.UpdateTag(ctx, db.UpdateTagParams{
		ID:   tagID,
		Name: normalizedName,
	})
	if err != nil {
		return nil, err
	}

	return s.mapTagToData(&tag), nil
}

func (s *Service) DeleteTag(ctx context.Context, tagID pgtype.UUID) error {
	// Check if tag is attached to any problems
	count, err := s.queries.CountTagUsage(ctx, tagID)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("TAG_IN_USE")
	}

	// Delete tag
	return s.queries.DeleteTag(ctx, tagID)
}

func (s *Service) AttachTagsToProblem(ctx context.Context, problemID pgtype.UUID, tagIDs []pgtype.UUID, orgID pgtype.UUID) error {
	// Verify all tags belong to the same organization
	tags, err := s.queries.GetTagsByIDs(ctx, tagIDs)
	if err != nil {
		return err
	}

	if len(tags) != len(tagIDs) {
		return errors.New("SOME_TAGS_NOT_FOUND")
	}

	for _, tag := range tags {
		if tag.OrganizationID.Bytes != orgID.Bytes {
			return errors.New("TAG_ORGANIZATION_MISMATCH")
		}
	}

	// Attach tags to problem
	for _, tagID := range tagIDs {
		err := s.queries.AddTagToProblem(ctx, db.AddTagToProblemParams{
			ProblemID: problemID,
			TagID:     tagID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DetachTagFromProblem(ctx context.Context, problemID, tagID pgtype.UUID) error {
	return s.queries.RemoveTagFromProblem(ctx, db.RemoveTagFromProblemParams{
		ProblemID: problemID,
		TagID:     tagID,
	})
}

// Resource operations

func (s *Service) AddResource(ctx context.Context, req CreateResourceRequest, problemID pgtype.UUID) (*ResourceData, error) {
	resource, err := s.queries.AddProblemResource(ctx, db.AddProblemResourceParams{
		ProblemID: problemID,
		Title:     pgtype.Text{String: req.Title, Valid: true},
		Url:       pgtype.Text{String: req.URL, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return s.mapResourceToData(&resource), nil
}

func (s *Service) ListResources(ctx context.Context, problemID pgtype.UUID) ([]ResourceData, error) {
	resources, err := s.queries.ListProblemResources(ctx, problemID)
	if err != nil {
		return nil, err
	}

	result := make([]ResourceData, len(resources))
	for i := range resources {
		result[i] = *s.mapResourceToData(&resources[i])
	}

	return result, nil
}

func (s *Service) GetResource(ctx context.Context, resourceID pgtype.UUID) (*db.ProblemResource, error) {
	resource, err := s.queries.GetProblemResource(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (s *Service) UpdateResource(ctx context.Context, req UpdateResourceRequest, resourceID pgtype.UUID) (*ResourceData, error) {
	// Check if at least one field is provided
	if req.Title == "" && req.URL == "" {
		return nil, errors.New("NO_FIELDS_PROVIDED")
	}

	// Build update params
	params := db.UpdateProblemResourceParams{
		ID: resourceID,
	}

	if req.Title != "" {
		params.Title = pgtype.Text{String: req.Title, Valid: true}
	}
	if req.URL != "" {
		params.Url = pgtype.Text{String: req.URL, Valid: true}
	}

	resource, err := s.queries.UpdateProblemResource(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.mapResourceToData(&resource), nil
}

func (s *Service) DeleteResource(ctx context.Context, resourceID pgtype.UUID) error {
	return s.queries.DeleteProblemResource(ctx, resourceID)
}

// Tag operations - to be implemented

// Resource operations - to be implemented

// Helper methods

func (s *Service) GetMember(ctx context.Context, orgID, userID pgtype.UUID) (*db.OrganizationMember, error) {
	member, err := s.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Service) mapProblemToData(problem *db.Problem) *ProblemData {
	return &ProblemData{
		ID:             problem.ID,
		OrganizationID: problem.OrganizationID,
		CreatedBy:      problem.CreatedBy,
		Title:          problem.Title,
		Description:    problem.Description.String,
		Difficulty:     string(problem.Difficulty),
		ExternalLink:   problem.ExternalLink.String,
		CreatedAt:      problem.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      problem.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		ArchivedAt:     formatTimestamp(problem.ArchivedAt),
	}
}

func (s *Service) mapTagToData(tag *db.Tag) *TagData {
	return &TagData{
		ID:             tag.ID,
		OrganizationID: tag.OrganizationID,
		Name:           tag.Name,
		CreatedAt:      tag.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) mapResourceToData(resource *db.ProblemResource) *ResourceData {
	return &ResourceData{
		ID:        resource.ID,
		ProblemID: resource.ProblemID,
		Title:     resource.Title.String,
		URL:       resource.Url.String,
		CreatedAt: resource.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func formatTimestamp(ts pgtype.Timestamptz) string {
	if ts.Valid {
		return ts.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return ""
}

// Super Admin operations

type ProblemWithOrgData struct {
	ProblemData
	OrganizationName string `json:"organization_name"`
	OrganizationSlug string `json:"organization_slug"`
}

func (s *Service) ListAllProblems(ctx context.Context, page, limit int) ([]ProblemWithOrgData, int, error) {
	// Calculate offset from page and limit
	offset := (page - 1) * limit

	// Get total count
	count, err := s.queries.CountAllProblems(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated problems
	problems, err := s.queries.ListAllProblems(ctx, db.ListAllProblemsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]ProblemWithOrgData, len(problems))
	for i, p := range problems {
		result[i] = ProblemWithOrgData{
			ProblemData: ProblemData{
				ID:             p.ID,
				OrganizationID: p.OrganizationID,
				Title:          p.Title,
				Description:    p.Description.String,
				Difficulty:     string(p.Difficulty),
				ExternalLink:   p.ExternalLink.String,
				CreatedBy:      p.CreatedBy,
				CreatedAt:      formatTimestamp(p.CreatedAt),
				UpdatedAt:      formatTimestamp(p.UpdatedAt),
			},
			OrganizationName: p.OrganizationName,
			OrganizationSlug: p.OrganizationSlug,
		}
	}

	return result, int(count), nil
}
