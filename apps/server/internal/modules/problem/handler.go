package problem

import (
	"net/http"

	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Problem handlers

// CreateProblem godoc
// @Summary Create a new problem
// @Description Create a new coding problem within an organization (mentor only)
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param body body CreateProblemRequest true "Problem details"
// @Success 201 {object} ProblemResponse "Problem created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - organization does not exist"
// @Router /v1/organizations/{orgId}/problems [post]
func (h *Handler) CreateProblem(c *echo.Context, body CreateProblemRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Prevent super_admin from creating content
	if claims.Role == "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_CANNOT_CREATE_CONTENT", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Create problem
	problem, err := h.service.CreateProblem((*c).Request().Context(), body, orgID, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ORGANIZATION_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "PROBLEM_CREATED", problem, nil)
}

// ListProblems godoc
// @Summary List problems
// @Description Get problems with filtering by difficulty, tags, and search query
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param difficulty query string false "Filter by difficulty (easy, medium, hard)"
// @Param tag_id query string false "Filter by tag ID (UUID)"
// @Param q query string false "Search by title"
// @Param sort_by query string false "Sort field (created_at, title, difficulty)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} ProblemListResponse "List of problems with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid organization ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/organizations/{orgId}/problems [get]
func (h *Handler) ListProblems(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// List problems
	problems, err := h.service.ListProblems((*c).Request().Context(), orgID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "PROBLEMS_RETRIEVED", problems, nil)
}

// GetProblem godoc
// @Summary Get problem by ID
// @Description Retrieve problem details including tags and resources
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Success 200 {object} ProblemResponse "Problem details"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId} [get]
func (h *Handler) GetProblem(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Get problem
	problem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	// Verify problem belongs to the organization
	if problem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "PROBLEM_RETRIEVED", problem, nil)
}

// UpdateProblem godoc
// @Summary Update problem details
// @Description Update problem information (mentor only)
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param body body UpdateProblemRequest true "Updated problem details"
// @Success 200 {object} ProblemResponse "Problem updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or no fields provided"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId} [patch]
func (h *Handler) UpdateProblem(c *echo.Context, body UpdateProblemRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Prevent super_admin from modifying content
	if claims.Role == "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_CANNOT_MODIFY_CONTENT", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Update problem
	problem, err := h.service.UpdateProblem((*c).Request().Context(), body, problemID)
	if err != nil {
		if err.Error() == "NO_FIELDS_PROVIDED" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "NO_FIELDS_PROVIDED", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "PROBLEM_UPDATED", problem, nil)
}

// DeleteProblem godoc
// @Summary Delete (archive) problem
// @Description Soft delete a problem using archived_at timestamp (mentor only)
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Success 200 {object} GenericResponse "Problem archived successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Failure 409 {object} map[string]any "Conflict - problem is referenced by assignments"
// @Router /v1/organizations/{orgId}/problems/{problemId} [delete]
func (h *Handler) DeleteProblem(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Prevent super_admin from deleting content
	if claims.Role == "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_CANNOT_DELETE_CONTENT", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Delete (archive) problem
	err = h.service.DeleteProblem((*c).Request().Context(), problemID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "PROBLEM_ARCHIVED", map[string]any{"message": "Problem archived successfully"}, nil)
}

// Tag handlers

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag for categorizing problems (mentor only)
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param body body CreateTagRequest true "Tag details"
// @Success 201 {object} TagResponse "Tag created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 409 {object} map[string]any "Conflict - tag name already exists in organization"
// @Router /v1/organizations/{orgId}/tags [post]
func (h *Handler) CreateTag(c *echo.Context, body CreateTagRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Create tag
	tag, err := h.service.CreateTag((*c).Request().Context(), body, orgID, userID)
	if err != nil {
		if err.Error() == "TAG_ALREADY_EXISTS" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "TAG_NAME_ALREADY_EXISTS", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "TAG_CREATED", tag, nil)
}

// ListTags godoc
// @Summary List tags
// @Description Get all tags for an organization with optional search
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param q query string false "Search by tag name"
// @Success 200 {object} TagListResponse "List of tags"
// @Failure 400 {object} map[string]any "Bad request - invalid organization ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Router /v1/organizations/{orgId}/tags [get]
func (h *Handler) ListTags(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Get search query parameter
	searchQuery := (*c).QueryParam("q")

	// List tags
	tags, err := h.service.ListTags((*c).Request().Context(), orgID, searchQuery)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "TAGS_RETRIEVED", tags, nil)
}

// UpdateTag godoc
// @Summary Update tag name
// @Description Update tag name (mentor only)
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param tagId path string true "Tag ID (UUID)"
// @Param body body UpdateTagRequest true "Updated tag details"
// @Success 200 {object} TagResponse "Tag updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - tag does not exist"
// @Failure 409 {object} map[string]any "Conflict - tag name already exists"
// @Router /v1/organizations/{orgId}/tags/{tagId} [patch]
func (h *Handler) UpdateTag(c *echo.Context, body UpdateTagRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	tagID, err := utils.StringToUUID((*c).Param("tagId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_TAG_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify tag exists and belongs to organization
	existingTag, err := h.service.GetTag((*c).Request().Context(), tagID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingTag.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
	}

	// Update tag
	tag, err := h.service.UpdateTag((*c).Request().Context(), body, tagID, orgID)
	if err != nil {
		if err.Error() == "TAG_NAME_ALREADY_EXISTS" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "TAG_NAME_ALREADY_EXISTS", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "TAG_UPDATED", tag, nil)
}

// DeleteTag godoc
// @Summary Delete tag
// @Description Delete a tag if not attached to any problems (mentor only)
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param tagId path string true "Tag ID (UUID)"
// @Success 200 {object} GenericResponse "Tag deleted successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - tag does not exist"
// @Failure 409 {object} map[string]any "Conflict - tag is attached to problems"
// @Router /v1/organizations/{orgId}/tags/{tagId} [delete]
func (h *Handler) DeleteTag(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	tagID, err := utils.StringToUUID((*c).Param("tagId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_TAG_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify tag exists and belongs to organization
	existingTag, err := h.service.GetTag((*c).Request().Context(), tagID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingTag.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
	}

	// Delete tag
	err = h.service.DeleteTag((*c).Request().Context(), tagID)
	if err != nil {
		if err.Error() == "TAG_IN_USE" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "TAG_IN_USE", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "TAG_DELETED", map[string]any{"message": "Tag deleted successfully"}, nil)
}

// AttachTagsToProblem godoc
// @Summary Attach tags to problem
// @Description Attach one or more tags to a problem (mentor only)
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param body body AttachTagsRequest true "Tag IDs to attach"
// @Success 200 {object} GenericResponse "Tags attached successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - problem or tags do not exist"
// @Failure 409 {object} map[string]any "Conflict - tags belong to different organization"
// @Router /v1/organizations/{orgId}/problems/{problemId}/tags [post]
func (h *Handler) AttachTagsToProblem(c *echo.Context, body AttachTagsRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Convert tag ID strings to UUIDs and deduplicate
	tagIDMap := make(map[string]pgtype.UUID)
	for _, tagIDStr := range body.TagIDs {
		tagID, err := utils.StringToUUID(tagIDStr)
		if err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_TAG_ID", nil, nil)
		}
		tagIDMap[tagIDStr] = tagID
	}

	// Convert map to slice
	tagIDs := make([]pgtype.UUID, 0, len(tagIDMap))
	for _, tagID := range tagIDMap {
		tagIDs = append(tagIDs, tagID)
	}

	// Attach tags to problem
	err = h.service.AttachTagsToProblem((*c).Request().Context(), problemID, tagIDs, orgID)
	if err != nil {
		if err.Error() == "SOME_TAGS_NOT_FOUND" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "SOME_TAGS_NOT_FOUND", nil, nil)
		}
		if err.Error() == "TAG_ORGANIZATION_MISMATCH" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "TAG_ORGANIZATION_MISMATCH", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "TAGS_ATTACHED", map[string]any{"message": "Tags attached successfully"}, nil)
}

// DetachTagFromProblem godoc
// @Summary Detach tag from problem
// @Description Remove a tag from a problem (mentor only)
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param tagId path string true "Tag ID (UUID)"
// @Success 200 {object} GenericResponse "Tag detached successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - problem or tag does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId}/tags/{tagId} [delete]
func (h *Handler) DetachTagFromProblem(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	tagID, err := utils.StringToUUID((*c).Param("tagId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_TAG_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Verify tag exists and belongs to organization
	existingTag, err := h.service.GetTag((*c).Request().Context(), tagID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingTag.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "TAG_NOT_FOUND", nil, nil)
	}

	// Detach tag from problem
	err = h.service.DetachTagFromProblem((*c).Request().Context(), problemID, tagID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "TAG_DETACHED", map[string]any{"message": "Tag detached successfully"}, nil)
}

// Resource handlers

// AddResource godoc
// @Summary Add resource to problem
// @Description Add a learning resource to a problem (mentor only)
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param body body CreateResourceRequest true "Resource details"
// @Success 201 {object} ResourceResponse "Resource added successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId}/resources [post]
func (h *Handler) AddResource(c *echo.Context, body CreateResourceRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Add resource
	resource, err := h.service.AddResource((*c).Request().Context(), body, problemID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "RESOURCE_ADDED", resource, nil)
}

// ListResources godoc
// @Summary List problem resources
// @Description Get all resources for a specific problem
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Success 200 {object} ResourceListResponse "List of resources"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId}/resources [get]
func (h *Handler) ListResources(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// List resources
	resources, err := h.service.ListResources((*c).Request().Context(), problemID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "RESOURCES_RETRIEVED", resources, nil)
}

// UpdateResource godoc
// @Summary Update resource
// @Description Update a problem resource (mentor only)
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param resourceId path string true "Resource ID (UUID)"
// @Param body body UpdateResourceRequest true "Updated resource details"
// @Success 200 {object} ResourceResponse "Resource updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or no fields provided"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - resource does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId}/resources/{resourceId} [patch]
func (h *Handler) UpdateResource(c *echo.Context, body UpdateResourceRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	resourceID, err := utils.StringToUUID((*c).Param("resourceId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_RESOURCE_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Verify resource exists and belongs to problem
	existingResource, err := h.service.GetResource((*c).Request().Context(), resourceID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "RESOURCE_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingResource.ProblemID.Bytes != problemID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "RESOURCE_NOT_FOUND", nil, nil)
	}

	// Update resource
	resource, err := h.service.UpdateResource((*c).Request().Context(), body, resourceID)
	if err != nil {
		if err.Error() == "NO_FIELDS_PROVIDED" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "NO_FIELDS_PROVIDED", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "RESOURCE_UPDATED", resource, nil)
}

// DeleteResource godoc
// @Summary Delete resource
// @Description Delete a problem resource (mentor only)
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param resourceId path string true "Resource ID (UUID)"
// @Success 200 {object} GenericResponse "Resource deleted successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - resource does not exist"
// @Router /v1/organizations/{orgId}/problems/{problemId}/resources/{resourceId} [delete]
func (h *Handler) DeleteResource(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	resourceID, err := utils.StringToUUID((*c).Param("resourceId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_RESOURCE_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is a member of the organization
	_, err = h.service.GetMember((*c).Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Verify problem exists and belongs to organization
	existingProblem, err := h.service.GetProblem((*c).Request().Context(), problemID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingProblem.OrganizationID.Bytes != orgID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
	}

	// Verify resource exists and belongs to problem
	existingResource, err := h.service.GetResource((*c).Request().Context(), resourceID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "RESOURCE_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	if existingResource.ProblemID.Bytes != problemID.Bytes {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "RESOURCE_NOT_FOUND", nil, nil)
	}

	// Delete resource
	err = h.service.DeleteResource((*c).Request().Context(), resourceID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "RESOURCE_DELETED", map[string]any{"message": "Resource deleted successfully"}, nil)
}

// Super Admin handlers

// ListAllProblems godoc
// @Summary List all problems (super admin only)
// @Description Retrieve all problems across all organizations with pagination
// @Tags Problems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} map[string]any "List of all problems with pagination"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - super admin role required"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/super-admin/problems [get]
func (h *Handler) ListAllProblems(c *echo.Context) error {
	// Validate super_admin role
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	if claims.Role != "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_ROLE_REQUIRED", nil, nil)
	}

	// Parse pagination parameters with defaults
	page := 1
	limit := 20

	if pageStr := (*c).QueryParam("page"); pageStr != "" {
		if p, err := utils.StringToInt(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := (*c).QueryParam("limit"); limitStr != "" {
		if l, err := utils.StringToInt(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, total, err := h.service.ListAllProblems(c.Request().Context(), page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
		"meta": map[string]any{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
