package mentorship

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	router := e.Group("/v1/mentorship")
	router.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	// User requests a role
	router.POST("/request-role", core.WithBody(handler.RequestRole))

	// Admin views requests
	router.GET("/requests", handler.ListRequests)

	// Admin approves/rejects requests
	router.PATCH("/requests/:id/status", core.WithBody(handler.UpdateStatus))
	router.DELETE("/requests/:id", handler.DeleteRequest)
}
