package routes

import (
	"net/http"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/container"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/organization"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Group, di *container.Container) {
	// health check api :
	e.GET("/health", healthCheck)

	auth.RegisterPublicRoutes(e, di.AuthHandler)
	auth.RegisterProtectedRoutes(e, di.AuthHandler, di.Config)

	organization.RegisterProtectedRoutes(e, di.OrganizationHandler, di.Config)
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
