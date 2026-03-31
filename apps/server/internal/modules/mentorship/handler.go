package mentorship

import (
	"log"
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RequestRole (Mentor/Mentee)
func (h *Handler) RequestRole(c *echo.Context, req RequestRoleRequest) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}
	userID := claims.UserID

	err := h.service.CreateRoleRequest(c.Request().Context(), userID, req.Role)
	if err != nil {
		log.Printf("Failed to create role request: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) ListRequests(c *echo.Context) error {
	// Need to handle admin check potentially, but skipping for MVP
	requests, err := h.service.ListPendingRequests(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *Handler) UpdateStatus(c *echo.Context, req UpdateStatusRequest) error {
	id := (*c).Param("id")
	
	err := h.service.UpdateStatus(c.Request().Context(), id, req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) DeleteRequest(c *echo.Context) error {
	id := (*c).Param("id")
	
	err := h.service.DeleteRequest(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
