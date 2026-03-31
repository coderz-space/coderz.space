package assignment

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	// Assignment Group routes
	groupRouter := e.Group("/v1/organizations/:orgId/bootcamps/:bootcampId/assignment-groups")
	groupRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	groupRouter.POST("", core.WithBody(handler.CreateAssignmentGroup))
	groupRouter.GET("", handler.ListAssignmentGroups)
	groupRouter.GET("/:groupId", handler.GetAssignmentGroup)
	groupRouter.PATCH("/:groupId", core.WithBody(handler.UpdateAssignmentGroup))
	groupRouter.DELETE("/:groupId", handler.DeleteAssignmentGroup)
	groupRouter.POST("/:groupId/problems", core.WithBody(handler.AddProblemsToGroup))
	groupRouter.PUT("/:groupId/problems", core.WithBody(handler.ReplaceGroupProblems))
	groupRouter.DELETE("/:groupId/problems/:problemId", handler.RemoveProblemFromGroup)

	// Assignment Instance routes
	assignmentRouter := e.Group("/v1/organizations/:orgId/bootcamps/:bootcampId/assignments")
	assignmentRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	assignmentRouter.POST("", core.WithBody(handler.CreateAssignment))
	assignmentRouter.GET("", handler.ListAssignments)
	assignmentRouter.GET("/:assignmentId", handler.GetAssignment)
	assignmentRouter.PATCH("/:assignmentId", core.WithBody(handler.UpdateAssignment))
	assignmentRouter.PATCH("/:assignmentId/deadline", core.WithBody(handler.UpdateAssignmentDeadline))
	assignmentRouter.PATCH("/:assignmentId/status", core.WithBody(handler.UpdateAssignmentStatus))

	// Assignment by enrollment routes
	enrollmentAssignmentRouter := e.Group("/v1/organizations/:orgId/bootcamps/:bootcampId/enrollments/:enrollmentId/assignments")
	enrollmentAssignmentRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	enrollmentAssignmentRouter.GET("", handler.ListAssignmentsByMentee)

	// Assignment Problem Progress routes
	problemProgressRouter := e.Group("/v1/organizations/:orgId/bootcamps/:bootcampId/assignments/:assignmentId/problems")
	problemProgressRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	problemProgressRouter.GET("", handler.ListAssignmentProblems)
	problemProgressRouter.GET("/:problemId", handler.GetAssignmentProblem)
	problemProgressRouter.PATCH("/:problemId", core.WithBody(handler.UpdateAssignmentProblemProgress))
}
