package assignment

import (
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
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

// Assignment Group Handlers

// CreateAssignmentGroup godoc
// @Summary Create a new assignment group
// @Description Create a reusable assignment template within a bootcamp (mentor only)
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param body body CreateAssignmentGroupRequest true "Assignment group details"
// @Success 201 {object} AssignmentGroupResponse "Assignment group created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups [post]
func (h *Handler) CreateAssignmentGroup(c *echo.Context, body CreateAssignmentGroupRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	createdBy, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	result, err := h.service.CreateAssignmentGroup((*c).Request().Context(), body, bootcampID, createdBy)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
		}
		if err.Error() == "BOOTCAMP_INACTIVE" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "BOOTCAMP_INACTIVE", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "ASSIGNMENT_GROUP_CREATED", result, nil)
}

// GetAssignmentGroup godoc
// @Summary Get assignment group details
// @Description Retrieve assignment group with associated problems
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param groupId path string true "Assignment Group ID (UUID)"
// @Success 200 {object} AssignmentGroupResponse "Assignment group details"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not a bootcamp member"
// @Failure 404 {object} map[string]any "Not found - assignment group does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId} [get]
func (h *Handler) GetAssignmentGroup(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	groupID, err := utils.StringToUUID((*c).Param("groupId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_GROUP_ID", nil, nil)
	}

	result, err := h.service.GetAssignmentGroup((*c).Request().Context(), groupID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_GROUP_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_GROUP_RETRIEVED", result, nil)
}

// UpdateAssignmentGroup godoc
// @Summary Update assignment group
// @Description Update assignment group details (title, description, deadline_days). Cannot change bootcamp_id. Does not affect existing assignment instances.
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param groupId path string true "Assignment Group ID (UUID)"
// @Param body body UpdateAssignmentGroupRequest true "Updated assignment group details"
// @Success 200 {object} AssignmentGroupResponse "Assignment group updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or no fields provided"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - assignment group does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId} [patch]
func (h *Handler) UpdateAssignmentGroup(c *echo.Context, body UpdateAssignmentGroupRequest) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	groupID, err := utils.StringToUUID((*c).Param("groupId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_GROUP_ID", nil, nil)
	}

	result, err := h.service.UpdateAssignmentGroup((*c).Request().Context(), groupID, body)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_GROUP_NOT_FOUND", nil, nil)
		}
		if err.Error() == "NO_FIELDS_PROVIDED" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "NO_FIELDS_PROVIDED", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_GROUP_UPDATED", result, nil)
}

// ListAssignmentGroups godoc
// @Summary List assignment groups
// @Description Get all assignment groups for a bootcamp with optional filtering and pagination
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param created_by query string false "Filter by creator user ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} AssignmentGroupListResponse "List of assignment groups with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid bootcamp ID or query parameters"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not a bootcamp member"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups [get]
func (h *Handler) ListAssignmentGroups(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	// Parse query parameters
	var createdBy *pgtype.UUID
	createdByStr := (*c).QueryParam("created_by")
	if createdByStr != "" {
		createdByUUID, err := utils.StringToUUID(createdByStr)
		if err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CREATED_BY_ID", nil, nil)
		}
		createdBy = &createdByUUID
	}

	// Parse pagination parameters
	page := 1
	if pageStr := (*c).QueryParam("page"); pageStr != "" {
		if p, err := utils.StringToInt(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := (*c).QueryParam("limit"); limitStr != "" {
		if l, err := utils.StringToInt(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	result, err := h.service.ListAssignmentGroups((*c).Request().Context(), bootcampID, createdBy, page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_GROUPS_RETRIEVED", result, nil)
}

// AddProblemsToGroup godoc
// @Summary Add problems to assignment group
// @Description Add or update problems in an assignment group with positions (mentor only)
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param groupId path string true "Assignment Group ID (UUID)"
// @Param body body AddProblemsToGroupRequest true "Problems to add with positions"
// @Success 200 {object} GenericResponse "Problems added successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - group or problem does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId}/problems [post]
func (h *Handler) AddProblemsToGroup(c *echo.Context, body AddProblemsToGroupRequest) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	groupID, err := utils.StringToUUID((*c).Param("groupId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_GROUP_ID", nil, nil)
	}

	err = h.service.AddProblemsToGroup((*c).Request().Context(), groupID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROBLEMS_ADDED_TO_GROUP", map[string]any{"message": "Problems added successfully"}, nil)
}

// RemoveProblemFromGroup godoc
// @Summary Remove problem from assignment group
// @Description Remove a problem from an assignment group (mentor only)
// @Tags Assignment Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param groupId path string true "Assignment Group ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Success 200 {object} GenericResponse "Problem removed successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - group or problem does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups/{groupId}/problems/{problemId} [delete]
func (h *Handler) RemoveProblemFromGroup(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	groupID, err := utils.StringToUUID((*c).Param("groupId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_GROUP_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	err = h.service.RemoveProblemFromGroup((*c).Request().Context(), groupID, problemID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROBLEM_REMOVED_FROM_GROUP", map[string]any{"message": "Problem removed successfully"}, nil)
}

// Assignment Instance Handlers

// CreateAssignment godoc
// @Summary Create assignment instance
// @Description Assign a problem set to a mentee with deadline (mentor only)
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param body body CreateAssignmentRequest true "Assignment details"
// @Success 201 {object} AssignmentResponse "Assignment created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - group or enrollment does not exist"
// @Failure 409 {object} map[string]any "Conflict - duplicate active assignment"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments [post]
func (h *Handler) CreateAssignment(c *echo.Context, body CreateAssignmentRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	assignedBy, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	result, err := h.service.CreateAssignment((*c).Request().Context(), body, assignedBy)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "ASSIGNMENT_CREATED", result, nil)
}

// GetAssignment godoc
// @Summary Get assignment details
// @Description Retrieve assignment with problem progress
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param assignmentId path string true "Assignment ID (UUID)"
// @Success 200 {object} AssignmentResponse "Assignment details with problems"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not authorized to view this assignment"
// @Failure 404 {object} map[string]any "Not found - assignment does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId} [get]
func (h *Handler) GetAssignment(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	assignmentID, err := utils.StringToUUID((*c).Param("assignmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_ID", nil, nil)
	}

	result, err := h.service.GetAssignment((*c).Request().Context(), assignmentID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_NOT_FOUND", nil, nil)
		}
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_RETRIEVED", result, nil)
}

// ListAssignmentsByMentee godoc
// @Summary List assignments for mentee
// @Description Get all assignments for a specific mentee enrollment
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param enrollmentId path string true "Bootcamp Enrollment ID (UUID)"
// @Success 200 {object} AssignmentListResponse "List of assignments"
// @Failure 400 {object} map[string]any "Bad request - invalid enrollment ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not authorized to view these assignments"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments/{enrollmentId}/assignments [get]
func (h *Handler) ListAssignmentsByMentee(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	enrollmentID, err := utils.StringToUUID((*c).Param("enrollmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ENROLLMENT_ID", nil, nil)
	}

	result, err := h.service.ListAssignmentsByMentee((*c).Request().Context(), enrollmentID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENTS_RETRIEVED", result, nil)
}

// UpdateAssignment godoc
// @Summary Update assignment
// @Description Update assignment status or deadline (mentor only)
// @Tags Assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param assignmentId path string true "Assignment ID (UUID)"
// @Param body body UpdateAssignmentRequest true "Updated assignment details"
// @Success 200 {object} AssignmentResponse "Assignment updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or no fields provided"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor role required"
// @Failure 404 {object} map[string]any "Not found - assignment does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId} [patch]
func (h *Handler) UpdateAssignment(c *echo.Context, body UpdateAssignmentRequest) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	assignmentID, err := utils.StringToUUID((*c).Param("assignmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_ID", nil, nil)
	}

	result, err := h.service.UpdateAssignment((*c).Request().Context(), assignmentID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_UPDATED", result, nil)
}

// Assignment Problem Progress Handlers

// UpdateAssignmentProblemProgress godoc
// @Summary Update problem progress
// @Description Update status, solution link, or notes for an assigned problem (mentee)
// @Tags Assignment Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param assignmentId path string true "Assignment ID (UUID)"
// @Param problemId path string true "Problem ID (UUID)"
// @Param body body UpdateAssignmentProblemRequest true "Progress update details"
// @Success 200 {object} AssignmentProblemResponse "Progress updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not authorized to update this problem"
// @Failure 404 {object} map[string]any "Not found - assignment problem does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId}/problems/{problemId} [patch]
func (h *Handler) UpdateAssignmentProblemProgress(c *echo.Context, body UpdateAssignmentProblemRequest) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	assignmentID, err := utils.StringToUUID((*c).Param("assignmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_ID", nil, nil)
	}

	problemID, err := utils.StringToUUID((*c).Param("problemId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
	}

	result, err := h.service.UpdateAssignmentProblemProgress((*c).Request().Context(), assignmentID, problemID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROGRESS_UPDATED", result, nil)
}

// ListAssignmentProblems godoc
// @Summary List assignment problems
// @Description Get all problems with progress for an assignment
// @Tags Assignment Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param assignmentId path string true "Assignment ID (UUID)"
// @Success 200 {object} AssignmentProblemListResponse "List of assignment problems with progress"
// @Failure 400 {object} map[string]any "Bad request - invalid assignment ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not authorized to view this assignment"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignments/{assignmentId}/problems [get]
func (h *Handler) ListAssignmentProblems(c *echo.Context) error {
	_, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	assignmentID, err := utils.StringToUUID((*c).Param("assignmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_ID", nil, nil)
	}

	result, err := h.service.ListAssignmentProblems((*c).Request().Context(), assignmentID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "ASSIGNMENT_PROBLEMS_RETRIEVED", result, nil)
}
