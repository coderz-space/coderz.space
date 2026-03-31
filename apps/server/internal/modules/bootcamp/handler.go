package bootcamp

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

// Bootcamp handlers

// CreateBootcamp godoc
// @Summary Create a new bootcamp
// @Description Create a new bootcamp within an organization (admin only)
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param body body CreateBootcampRequest true "Bootcamp details"
// @Success 201 {object} BootcampResponse "Bootcamp created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid date range"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 404 {object} map[string]any "Not found - organization does not exist"
// @Failure 409 {object} map[string]any "Conflict - organization not approved"
// @Router /v1/organizations/{orgId}/bootcamps [post]
func (h *Handler) CreateBootcamp(c *echo.Context) error {
	var body CreateBootcampRequest
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

	// Prevent super_admin from creating bootcamps
	if claims.Role == "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_CANNOT_CREATE_CONTENT", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member ID for created_by
	memberID, err := h.service.GetMemberID(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	data, err := h.service.CreateBootcamp(c.Request().Context(), orgID, body, memberID)
	if err != nil {
		if err.Error() == "ORGANIZATION_NOT_FOUND" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ORGANIZATION_NOT_FOUND", nil, nil)
		}
		if err.Error() == "ORGANIZATION_NOT_APPROVED" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "ORGANIZATION_NOT_APPROVED", nil, nil)
		}
		if err.Error() == "INVALID_DATE_RANGE" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "START_DATE_MUST_BE_BEFORE_END_DATE", nil, nil)
		}
		if err.Error() == "INVALID_START_DATE" || err.Error() == "INVALID_END_DATE" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusCreated, BootcampResponse{
		Success: true,
		Data:    *data,
	})
}

// GetBootcamp godoc
// @Summary Get bootcamp by ID
// @Description Retrieve bootcamp details by ID with role-based access control
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Success 200 {object} BootcampResponse "Bootcamp details"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist or not enrolled"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId} [get]
func (h *Handler) GetBootcamp(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member to determine role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Fetch bootcamp details
	data, err := h.service.GetBootcampByID(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	// Validate bootcamp belongs to the organization (cross-org access check)
	if data.OrganizationID != orgID {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	// Role-based access validation
	if member.Role == "mentee" {
		// Mentees can only access bootcamps where they are enrolled
		_, err := h.service.GetEnrollmentByMember(c.Request().Context(), bootcampID, member.ID)
		if err != nil {
			// Return 404 if mentee is not enrolled (not 403 to avoid information disclosure)
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
		}
	}
	// Admins and mentors can access any bootcamp in their organization

	return c.JSON(http.StatusOK, BootcampResponse{
		Success: true,
		Data:    *data,
	})
}

// ListBootcamps godoc
// @Summary List bootcamps
// @Description Get bootcamps with role-based filtering (mentees see only enrolled bootcamps)
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param is_active query boolean false "Filter by active status"
// @Success 200 {object} BootcampListResponse "List of bootcamps with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid organization ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not an organization member"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/organizations/{orgId}/bootcamps [get]
func (h *Handler) ListBootcamps(c *echo.Context) error {
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
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
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

	// Parse is_active filter
	var isActive *bool
	if isActiveStr := (*c).QueryParam("is_active"); isActiveStr != "" {
		switch isActiveStr {
		case "true":
			val := true
			isActive = &val
		case "false":
			val := false
			isActive = &val
		}
	}

	// Get the organization member to determine role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Determine filtering based on role
	var memberID *pgtype.UUID
	if member.Role == "mentee" {
		// Mentees only see bootcamps where they are enrolled
		memberID = &member.ID
	}
	// Admins and mentors see all bootcamps in the organization (memberID = nil)

	data, total, err := h.service.ListBootcampsWithFilters(c.Request().Context(), orgID, memberID, isActive, page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, BootcampListResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

// UpdateBootcamp godoc
// @Summary Update bootcamp details
// @Description Update bootcamp information (admin only)
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param body body UpdateBootcampRequest true "Updated bootcamp details"
// @Success 200 {object} BootcampResponse "Bootcamp updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or no fields provided"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - admin role required"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId} [patch]
func (h *Handler) UpdateBootcamp(c *echo.Context) error {
	var body UpdateBootcampRequest
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

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member to verify admin role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Validate admin role authorization
	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	// Verify bootcamp belongs to the organization
	bootcamp, err := h.service.GetBootcampByID(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	if bootcamp.OrganizationID != orgID {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	data, err := h.service.UpdateBootcamp(c.Request().Context(), bootcampID, body)
	if err != nil {
		if err.Error() == "NO_FIELDS_PROVIDED" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "NO_FIELDS_PROVIDED", nil, nil)
		}
		if err.Error() == "INVALID_DATE_RANGE" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "START_DATE_MUST_BE_BEFORE_END_DATE", nil, nil)
		}
		if err.Error() == "INVALID_START_DATE" || err.Error() == "INVALID_END_DATE" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, BootcampResponse{
		Success: true,
		Data:    *data,
	})
}

// DeactivateBootcamp godoc
// @Summary Deactivate bootcamp
// @Description Set bootcamp is_active to false (admin only)
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Success 200 {object} GenericResponse "Bootcamp deactivated successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - admin role required"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/deactivate [post]
func (h *Handler) DeactivateBootcamp(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member to verify admin role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Validate admin role authorization
	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	// Verify bootcamp belongs to the organization
	bootcamp, err := h.service.GetBootcampByID(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	if bootcamp.OrganizationID != orgID {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	err = h.service.DeactivateBootcamp(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}

// Enrollment handlers

// EnrollMember godoc
// @Summary Enroll member in bootcamp
// @Description Enroll an organization member into a bootcamp with specified role (admin only)
// @Tags Bootcamp Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param body body EnrollMemberRequest true "Enrollment details"
// @Success 201 {object} EnrollmentResponse "Member enrolled successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - admin role required"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist"
// @Failure 409 {object} map[string]any "Conflict - bootcamp inactive or cross-org violation"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments [post]
func (h *Handler) EnrollMember(c *echo.Context) error {
	var body EnrollMemberRequest
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

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member to verify admin role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Validate admin role authorization
	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	data, err := h.service.EnrollMember(c.Request().Context(), orgID, bootcampID, body)
	if err != nil {
		if err.Error() == "BOOTCAMP_INACTIVE" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "BOOTCAMP_INACTIVE", nil, nil)
		}
		if err.Error() == "BOOTCAMP_NOT_FOUND" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
		}
		if err.Error() == "CROSS_ORG_VIOLATION" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "CROSS_ORG_VIOLATION", nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusCreated, EnrollmentResponse{
		Success: true,
		Data:    *data,
	})
}

// ListEnrollments godoc
// @Summary List bootcamp enrollments
// @Description Get all enrollments for a bootcamp
// @Tags Bootcamp Enrollments
// @Accept json
// @Produce json
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Success 200 {object} EnrollmentListResponse "List of enrollments"
// @Failure 400 {object} map[string]any "Bad request - invalid bootcamp ID"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/bootcamps/{bootcampId}/enrollments [get]
func (h *Handler) ListEnrollments(c *echo.Context) error {
	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	data, err := h.service.ListEnrollments(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, EnrollmentListResponse{
		Success: true,
		Data:    data,
	})
}

// UpdateEnrollmentRole godoc
// @Summary Update enrollment role
// @Description Update the role of a bootcamp enrollment (admin only)
// @Tags Bootcamp Enrollments
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param enrollmentId path string true "Enrollment ID (UUID)"
// @Param body body UpdateEnrollmentRoleRequest true "New role"
// @Success 200 {object} EnrollmentResponse "Enrollment role updated successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments/{enrollmentId} [patch]
func (h *Handler) UpdateEnrollmentRole(c *echo.Context) error {
	var body UpdateEnrollmentRoleRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}

	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	enrollmentID, err := utils.StringToUUID(c.Param("enrollmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ENROLLMENT_ID", nil, nil)
	}

	data, err := h.service.UpdateEnrollmentRole(c.Request().Context(), enrollmentID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, EnrollmentResponse{
		Success: true,
		Data:    *data,
	})
}

// RemoveEnrollment godoc
// @Summary Remove enrollment
// @Description Remove a member's enrollment from a bootcamp (admin only)
// @Tags Bootcamp Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organization ID (UUID)"
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param enrollmentId path string true "Enrollment ID (UUID)"
// @Success 200 {object} GenericResponse "Enrollment removed successfully"
// @Failure 400 {object} map[string]any "Bad request - invalid enrollment ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - admin role required"
// @Failure 404 {object} map[string]any "Not found - enrollment does not exist"
// @Router /v1/organizations/{orgId}/bootcamps/{bootcampId}/enrollments/{enrollmentId} [delete]
func (h *Handler) RemoveEnrollment(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	enrollmentID, err := utils.StringToUUID((*c).Param("enrollmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ENROLLMENT_ID", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get the organization member to verify admin role
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	// Validate admin role authorization
	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	// Verify enrollment exists and belongs to the bootcamp
	enrollment, err := h.service.GetEnrollment(c.Request().Context(), enrollmentID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ENROLLMENT_NOT_FOUND", nil, nil)
	}

	if enrollment.BootcampID != bootcampID {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ENROLLMENT_NOT_FOUND", nil, nil)
	}

	// Verify bootcamp belongs to the organization
	bootcamp, err := h.service.GetBootcampByID(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	if bootcamp.OrganizationID != orgID {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	err = h.service.RemoveEnrollment(c.Request().Context(), enrollmentID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}

// Super Admin handlers

// ListAllBootcamps godoc
// @Summary List all bootcamps (super admin only)
// @Description Retrieve all bootcamps across all organizations with pagination
// @Tags Bootcamps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} map[string]any "List of all bootcamps with pagination"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - super admin role required"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /v1/super-admin/bootcamps [get]
func (h *Handler) ListAllBootcamps(c *echo.Context) error {
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

	data, total, err := h.service.ListAllBootcamps(c.Request().Context(), page, limit)
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
