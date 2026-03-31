package routes

import (
	"net/http"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/container"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/analytics"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/assignment"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/bootcamp"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/organization"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/problem"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/progress"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Group, di *container.Container) {
	// health check api :
	e.GET("/health", healthCheck)

	// Auth module routes (public and protected)
	auth.RegisterPublicRoutes(e, di.AuthHandler)
	auth.RegisterProtectedRoutes(e, di.AuthHandler, di.Config)

	// Organization module routes
	organization.RegisterProtectedRoutes(e, di.OrganizationHandler, di.Config)

	// Bootcamp module routes
	bootcamp.RegisterProtectedRoutes(e, di.BootcampHandler, di.Config)

	// Problem module routes
	problem.RegisterProtectedRoutes(e, di.ProblemHandler, di.Config)

	// Assignment module routes
	assignment.RegisterProtectedRoutes(e, di.AssignmentHandler, di.Config)

	// Progress (doubts) module routes
	progress.RegisterProtectedRoutes(e, di.ProgressHandler, di.Config)

	// Analytics module routes
	analytics.RegisterProtectedRoutes(e, di.AnalyticsHandler, di.Config)
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
