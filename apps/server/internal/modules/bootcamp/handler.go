package bootcamp

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

// Bootcamp handlers

func (h *Handler) CreateBootcamp(c *echo.Context, body CreateBootcampRequest) error {
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

func (h *Handler) UpdateBootcamp(c *echo.Context, body UpdateBootcampRequest) error {
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

func (h *Handler) EnrollMember(c *echo.Context, body EnrollMemberRequest) error {
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

func (h *Handler) UpdateEnrollmentRole(c *echo.Context, body UpdateEnrollmentRoleRequest) error {
	enrollmentID, err := utils.StringToUUID((*c).Param("enrollmentId"))
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

func (h *Handler) RemoveEnrollment(c *echo.Context) error {
	enrollmentID, err := utils.StringToUUID((*c).Param("enrollmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ENROLLMENT_ID", nil, nil)
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
