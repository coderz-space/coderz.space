package profile

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	router := e.Group("/v1")
	router.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	router.GET("/mentees/:id/profile", handler.GetMenteeProfile)
	router.GET("/leaderboard", handler.GetLeaderboard)
}
