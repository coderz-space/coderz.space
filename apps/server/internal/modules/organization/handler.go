package organization

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

// Organization handlers

func (h *Handler) CreateOrganization(c *echo.Context, body CreateOrganizationRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	data, err := h.service.CreateOrganization(c.Request().Context(), body, userID)
	if err != nil {
		if err.Error() == "SLUG_ALREADY_EXISTS" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "SLUG_ALREADY_EXISTS", nil, nil)
		}
		if err.Error() == "INVALID_SLUG_FORMAT" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_SLUG_FORMAT", nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusCreated, OrganizationResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) GetOrganization(c *echo.Context) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	data, err := h.service.GetOrganizationByID(c.Request().Context(), orgID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ORGANIZATION_NOT_FOUND", nil, nil)
	}

	return c.JSON(http.StatusOK, OrganizationResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) ListOrganizations(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
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

	data, total, err := h.service.ListUserOrganizations(c.Request().Context(), userID, page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, OrganizationListResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

func (h *Handler) UpdateOrganization(c *echo.Context, body UpdateOrganizationRequest) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	// Get authenticated user
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Check if user is an admin of the organization
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	data, err := h.service.UpdateOrganization(c.Request().Context(), orgID, body)
	if err != nil {
		if err.Error() == "NO_FIELDS_PROVIDED" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "NO_FIELDS_PROVIDED", nil, nil)
		}
		if err.Error() == "SLUG_ALREADY_EXISTS" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "SLUG_ALREADY_EXISTS", nil, nil)
		}
		if err.Error() == "INVALID_SLUG_FORMAT" {
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_SLUG_FORMAT", nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, OrganizationResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) ApproveOrganization(c *echo.Context) error {
	// Validate super_admin role
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	if claims.Role != "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_ROLE_REQUIRED", nil, nil)
	}

	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	data, err := h.service.ApproveOrganization(c.Request().Context(), orgID)
	if err != nil {
		if err.Error() == "ORGANIZATION_NOT_FOUND" {
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ORGANIZATION_NOT_FOUND", nil, nil)
		}
		if err.Error() == "ORGANIZATION_NOT_PENDING" {
			return response.NewResponse(c, http.StatusConflict, "CONFLICT", "ORGANIZATION_NOT_PENDING", nil, nil)
		}
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, OrganizationResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) GetPendingOrganizations(c *echo.Context) error {
	// Validate super_admin role
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	if claims.Role != "super_admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "SUPER_ADMIN_ROLE_REQUIRED", nil, nil)
	}

	data, err := h.service.GetPendingOrganizations(c.Request().Context())
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, OrganizationListResponse{
		Success: true,
		Data:    data,
	})
}

// Member handlers

func (h *Handler) AddMember(c *echo.Context, body AddMemberRequest) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	// Get authenticated user
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Check if user is an admin of the organization
	member, err := h.service.GetMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ORGANIZATION_MEMBER", nil, nil)
	}

	if member.Role != "admin" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ADMIN_ROLE_REQUIRED", nil, nil)
	}

	data, err := h.service.AddMember(c.Request().Context(), orgID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusCreated, MemberResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) ListMembers(c *echo.Context) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
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

	data, total, err := h.service.ListMembers(c.Request().Context(), orgID, page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, MemberListResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

func (h *Handler) UpdateMemberRole(c *echo.Context, body UpdateMemberRoleRequest) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID((*c).Param("userId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	data, err := h.service.UpdateMemberRole(c.Request().Context(), orgID, userID, body)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, MemberResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) RemoveMember(c *echo.Context) error {
	orgID, err := utils.StringToUUID((*c).Param("orgId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ORGANIZATION_ID", nil, nil)
	}

	userID, err := utils.StringToUUID((*c).Param("userId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_USER_ID", nil, nil)
	}

	err = h.service.RemoveMember(c.Request().Context(), orgID, userID)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}
