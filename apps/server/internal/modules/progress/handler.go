package progress

import (
	"net/http"

	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/coderz-space/coderz.space/internal/common/validator"
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

// CreateDoubt godoc
// @Summary Create a new doubt
// @Description Create a doubt for an assignment problem (mentee only). Rate limited to prevent spam.
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateDoubtRequest true "Doubt details"
// @Success 201 {object} DoubtResponse "Doubt created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid assignment problem ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not a mentee or problem not assigned to you"
// @Failure 404 {object} map[string]any "Not found - assignment problem does not exist"
// @Failure 429 {object} map[string]any "Too many requests - rate limit exceeded"
// @Router /v1/doubts [post]
func (h *Handler) CreateDoubt(c *echo.Context) error {
	var body CreateDoubtRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}

	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Parse assignment_problem_id to get bootcamp context
	assignmentProblemID, err := utils.StringToUUID(body.AssignmentProblemID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_PROBLEM_ID", nil, nil)
	}

	// Get assignment problem details to find bootcamp
	apDetails, err := h.service.queries.GetAssignmentProblemDetails(c.Request().Context(), assignmentProblemID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_PROBLEM_NOT_FOUND", nil, nil)
	}

	// Get member ID for the user in this bootcamp
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	memberID, err := h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, apDetails.BootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Create doubt
	doubt, err := h.service.CreateDoubt(c.Request().Context(), body, memberID)
	if err != nil {
		switch err.Error() {
		case "INVALID_ASSIGNMENT_PROBLEM_ID":
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ASSIGNMENT_PROBLEM_ID", nil, nil)
		case "ASSIGNMENT_PROBLEM_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_PROBLEM_NOT_FOUND", nil, nil)
		case "ASSIGNMENT_PROBLEM_NOT_OWNED":
			return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ASSIGNMENT_PROBLEM_NOT_OWNED", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusCreated, "SUCCESS", "DOUBT_CREATED", doubt, nil)
}

// ListDoubts godoc
// @Summary List doubts
// @Description List doubts with filtering and cursor-based pagination. Mentees see only their own doubts, mentors/admins see all organization doubts.
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId query string false "Filter by bootcamp ID (UUID) - required for mentors/admins"
// @Param assignmentProblemId query string false "Filter by assignment problem ID (UUID)"
// @Param resolved query boolean false "Filter by resolved status"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Success 200 {object} DoubtListResponse "List of doubts with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid query parameters"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - insufficient permissions"
// @Router /v1/doubts [get]
func (h *Handler) ListDoubts(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Parse query parameters
	bootcampIDStr := (*c).QueryParam("bootcampId")
	assignmentProblemIDStr := (*c).QueryParam("assignmentProblemId")
	resolvedStr := (*c).QueryParam("resolved")
	cursor := (*c).QueryParam("cursor")
	limitStr := (*c).QueryParam("limit")

	// Parse limit
	limit := ParseLimit(limitStr, 20, 100)

	// Build filters
	filters := make(map[string]string)
	if assignmentProblemIDStr != "" {
		filters["assignment_problem_id"] = assignmentProblemIDStr
	}
	if resolvedStr != "" {
		filters["resolved"] = resolvedStr
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Determine user role and get member ID
	userRole := claims.Role
	var bootcampID pgtype.UUID
	var memberID pgtype.UUID

	if bootcampIDStr == "" {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "BOOTCAMP_ID_REQUIRED", nil, nil)
	}
	bootcampID, err = utils.StringToUUID(bootcampIDStr)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	memberID, err = h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// List doubts
	doubts, pagination, err := h.service.ListDoubts(c.Request().Context(), bootcampID, filters, limit, cursor, userRole, memberID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "DOUBTS_RETRIEVED", DoubtListResponse{
		Success: true,
		Data:    doubts,
		Meta:    pagination,
	}, nil)
}

// GetDoubt godoc
// @Summary Get doubt details
// @Description Retrieve full details of a specific doubt. Mentees can only view their own doubts, mentors/admins can view all organization doubts.
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param doubtId path string true "Doubt ID (UUID)"
// @Success 200 {object} DoubtResponse "Doubt details"
// @Failure 400 {object} map[string]any "Bad request - invalid doubt ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - access denied"
// @Failure 404 {object} map[string]any "Not found - doubt does not exist"
// @Router /v1/doubts/{doubtId} [get]
func (h *Handler) GetDoubt(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	doubtID, err := utils.StringToUUID((*c).Param("doubtId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_DOUBT_ID", nil, nil)
	}

	// Get user ID and role
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// First, get the doubt to find which bootcamp it belongs to
	doubt, err := h.service.queries.GetDoubt(c.Request().Context(), doubtID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "DOUBT_NOT_FOUND", nil, nil)
	}

	// Get assignment problem details to find bootcamp
	apDetails, err := h.service.queries.GetAssignmentProblemDetails(c.Request().Context(), doubt.AssignmentProblemID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_PROBLEM_NOT_FOUND", nil, nil)
	}

	// Get member ID for access control
	memberID, err := h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, apDetails.BootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Get doubt with access control
	doubtData, err := h.service.GetDoubt(c.Request().Context(), doubtID, claims.Role, memberID)
	if err != nil {
		switch err.Error() {
		case "DOUBT_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "DOUBT_NOT_FOUND", nil, nil)
		case "ACCESS_DENIED":
			return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ACCESS_DENIED", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "DOUBT_RETRIEVED", doubtData, nil)
}

// ResolveDoubt godoc
// @Summary Resolve a doubt
// @Description Mark a doubt as resolved by a mentor/admin with optional resolution note. Idempotent operation.
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param doubtId path string true "Doubt ID (UUID)"
// @Param body body ResolveDoubtRequest true "Resolution details"
// @Success 200 {object} DoubtResponse "Doubt resolved successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid doubt ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - only mentors/admins can resolve doubts"
// @Failure 404 {object} map[string]any "Not found - doubt does not exist"
// @Router /v1/doubts/{doubtId}/resolve [patch]
func (h *Handler) ResolveDoubt(c *echo.Context) error {
	var body ResolveDoubtRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}

	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentor/admin
	if claims.Role == "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ONLY_MENTORS_ADMINS_CAN_RESOLVE", nil, nil)
	}

	doubtID, err := utils.StringToUUID((*c).Param("doubtId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_DOUBT_ID", nil, nil)
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get doubt to find bootcamp
	doubt, err := h.service.queries.GetDoubt(c.Request().Context(), doubtID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "DOUBT_NOT_FOUND", nil, nil)
	}

	// Get assignment problem details to find bootcamp
	apDetails, err := h.service.queries.GetAssignmentProblemDetails(c.Request().Context(), doubt.AssignmentProblemID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ASSIGNMENT_PROBLEM_NOT_FOUND", nil, nil)
	}

	// Get resolver's member ID
	resolverMemberID, err := h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, apDetails.BootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Resolve doubt
	resolvedDoubt, err := h.service.ResolveDoubt(c.Request().Context(), doubtID, resolverMemberID, body.ResolutionNote)
	if err != nil {
		switch err.Error() {
		case "DOUBT_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "DOUBT_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "DOUBT_RESOLVED", resolvedDoubt, nil)
}

// DeleteDoubt godoc
// @Summary Delete a doubt
// @Description Permanently delete a doubt (mentor/admin only). Mentees cannot delete doubts for audit purposes.
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param doubtId path string true "Doubt ID (UUID)"
// @Success 200 {object} GenericResponse "Doubt deleted successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid doubt ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - only mentors/admins can delete doubts"
// @Failure 404 {object} map[string]any "Not found - doubt does not exist"
// @Router /v1/doubts/{doubtId} [delete]
func (h *Handler) DeleteDoubt(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentor/admin
	if claims.Role == "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "MENTEES_CANNOT_DELETE_DOUBTS", nil, nil)
	}

	doubtID, err := utils.StringToUUID((*c).Param("doubtId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_DOUBT_ID", nil, nil)
	}

	// Delete doubt
	err = h.service.DeleteDoubt(c.Request().Context(), doubtID, claims.Role)
	if err != nil {
		switch err.Error() {
		case "DOUBT_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "DOUBT_NOT_FOUND", nil, nil)
		case "MENTEES_CANNOT_DELETE_DOUBTS":
			return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "MENTEES_CANNOT_DELETE_DOUBTS", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "DOUBT_DELETED", map[string]any{
		"message": "Doubt deleted successfully",
	}, nil)
}

// GetMyDoubts godoc
// @Summary Get my doubts
// @Description Retrieve all doubts raised by the authenticated mentee with cursor-based pagination
// @Tags Doubts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId query string true "Bootcamp ID (UUID)"
// @Param resolved query boolean false "Filter by resolved status"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Success 200 {object} DoubtListResponse "List of my doubts with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid query parameters"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - only mentees can access this endpoint"
// @Router /v1/doubts/me [get]
func (h *Handler) GetMyDoubts(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Parse query parameters
	bootcampIDStr := (*c).QueryParam("bootcampId")
	if bootcampIDStr == "" {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "BOOTCAMP_ID_REQUIRED", nil, nil)
	}

	bootcampID, err := utils.StringToUUID(bootcampIDStr)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	resolvedStr := (*c).QueryParam("resolved")
	cursor := (*c).QueryParam("cursor")
	limitStr := (*c).QueryParam("limit")

	// Parse limit
	limit := ParseLimit(limitStr, 20, 100)

	// Build filters
	filters := make(map[string]string)
	if resolvedStr != "" {
		filters["resolved"] = resolvedStr
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get member ID
	memberID, err := h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// List doubts (force mentee role to filter by raised_by)
	doubts, pagination, err := h.service.ListDoubts(c.Request().Context(), bootcampID, filters, limit, cursor, "mentee", memberID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "DOUBTS_RETRIEVED", DoubtListResponse{
		Success: true,
		Data:    doubts,
		Meta:    pagination,
	}, nil)
}
