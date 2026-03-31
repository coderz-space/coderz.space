package bootcamp

import (
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
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
	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	data, err := h.service.GetBootcampByID(c.Request().Context(), bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
	}

	return c.JSON(http.StatusOK, BootcampResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) ListBootcamps(c *echo.Context) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	data, err := h.service.ListBootcampsByOrg(c.Request().Context(), orgID)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, BootcampListResponse{
		Success: true,
		Data:    data,
	})
}

func (h *Handler) UpdateBootcamp(c *echo.Context, body UpdateBootcampRequest) error {
	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
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
	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
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
	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	data, err := h.service.EnrollMember(c.Request().Context(), bootcampID, body)
	if err != nil {
		if err.Error() == "BOOTCAMP_INACTIVE" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "BOOTCAMP_INACTIVE", nil, nil)
		}
		if err.Error() == "BOOTCAMP_NOT_FOUND" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
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
