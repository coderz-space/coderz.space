package profile

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMenteeProfile(c *echo.Context) error {
	id := (*c).Param("id")
	res, err := h.service.GetMenteeProfile(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetLeaderboard(c *echo.Context) error {
	res, err := h.service.GetLeaderboard(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
