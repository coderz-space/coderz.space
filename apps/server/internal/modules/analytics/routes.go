package analytics

import (
	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

// RegisterProtectedRoutes registers all analytics module routes
func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	// Leaderboard routes
	leaderboardRouter := e.Group("/v1/bootcamps/:bootcampId/leaderboard")
	leaderboardRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	leaderboardRouter.GET("", handler.GetBootcampLeaderboard)            // Get bootcamp leaderboard
	leaderboardRouter.GET("/:enrollmentId", handler.GetLeaderboardEntry) // Get specific leaderboard entry

	// Poll routes
	pollRouter := e.Group("/v1/bootcamps/:bootcampId/polls")
	pollRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	pollRouter.POST("", handler.CreatePoll)                    // Create poll (mentor/admin only)
	pollRouter.GET("", handler.ListPolls)                      // List polls
	pollRouter.GET("/:pollId", handler.GetPoll)                // Get poll details
	pollRouter.PUT("/:pollId/vote", handler.VotePoll)          // Vote on poll (mentee only)
	pollRouter.GET("/:pollId/results", handler.GetPollResults) // Get poll results (mentor/admin only)
	pollRouter.GET("/:pollId/votes", handler.GetPollVotes)     // Get individual votes (mentor/admin only)
}
