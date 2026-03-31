package tasks

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	router := e.Group("/v1")
	router.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	router.GET("/mentees/:id/questions", handler.ListMenteeQuestions)
	router.POST("/mentees/:id/questions", core.WithBody(handler.AssignQuestionToMentee))
	router.GET("/mentees/:id/questions/:questionId", handler.GetMenteeQuestion)
	
	router.PATCH("/mentees/:id/questions/:questionId/progress", core.WithBody(handler.UpdateQuestionProgress))
	router.PATCH("/mentees/:id/questions/:questionId/details", core.WithBody(handler.UpdateQuestionDetails))
}
