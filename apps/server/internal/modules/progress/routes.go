package progress

import (
	"time"

	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/middleware/ratelimit"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

// RegisterProtectedRoutes registers all progress (doubts) module routes
func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	doubtRouter := e.Group("/v1/doubts")
	doubtRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	// Doubt management endpoints with rate limiting for creation
	// Rate limit: 10 requests per minute per user for doubt creation
	doubtRouter.POST("", handler.CreateDoubt, ratelimit.RateLimitMiddleware(10, 10, time.Minute)) // Create doubt (mentee only)
	doubtRouter.GET("", handler.ListDoubts)                                                       // List doubts (role-based filtering)
	doubtRouter.GET("/me", handler.GetMyDoubts)                                                   // Get my doubts (mentee only)
	doubtRouter.GET("/:doubtId", handler.GetDoubt)                                                // Get doubt details
	doubtRouter.PATCH("/:doubtId/resolve", handler.ResolveDoubt)                                  // Resolve doubt (mentor/admin only)
	doubtRouter.DELETE("/:doubtId", handler.DeleteDoubt)                                          // Delete doubt (mentor/admin only)
}
