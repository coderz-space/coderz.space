package app

import (
	authmw "github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	appRouter := e.Group("/v1/app")
	appRouter.POST("/auth/mentee-signup", handler.MenteeSignup)
}

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	appRouter := e.Group("/v1/app")
	appRouter.Use(authmw.AuthMiddleware(config.JWTSecret, config.JWTExpires))

	appRouter.GET("/context", handler.GetContext)
	appRouter.GET("/mentor/mentee-requests", handler.ListMenteeRequests)
	appRouter.PATCH("/mentor/mentee-requests/:requestId", handler.ReviewMenteeRequest)
	appRouter.GET("/sheets", handler.ListSheets)
	appRouter.GET("/mentor/day-assignments/:day", handler.GetDayAssignments)
	appRouter.PUT("/mentor/day-assignments/:day", handler.UpdateDayAssignments)
	appRouter.POST("/mentor/assignments", handler.CreateAssignments)
	appRouter.GET("/mentees/:username/questions", handler.ListMenteeQuestions)
	appRouter.GET("/mentees/:username/questions/:assignmentProblemId", handler.GetMenteeQuestion)
	appRouter.PATCH("/mentees/:username/questions/:assignmentProblemId", handler.UpdateMenteeQuestion)
	appRouter.GET("/mentees/:username/profile", handler.GetMenteeProfile)
	appRouter.GET("/me/profile", handler.GetMyProfile)
	appRouter.PATCH("/me/profile", handler.UpdateMyProfile)
	appRouter.PATCH("/me/password", handler.UpdateMyPassword)
	appRouter.GET("/leaderboard", handler.GetLeaderboard)
}
