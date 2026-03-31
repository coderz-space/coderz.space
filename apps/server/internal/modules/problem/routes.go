package problem

import (
	"github.com/coderz-space/coderz.space/internal/common/core"
	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	// Problem routes
	problemRouter := e.Group("/v1/organizations/:orgId/problems")
	problemRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	problemRouter.POST("", core.WithBody(handler.CreateProblem))
	problemRouter.GET("", handler.ListProblems)
	problemRouter.GET("/:problemId", handler.GetProblem)
	problemRouter.PATCH("/:problemId", core.WithBody(handler.UpdateProblem))
	problemRouter.DELETE("/:problemId", handler.DeleteProblem)

	// Tag routes
	tagRouter := e.Group("/v1/organizations/:orgId/tags")
	tagRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	tagRouter.POST("", core.WithBody(handler.CreateTag))
	tagRouter.GET("", handler.ListTags)
	tagRouter.PATCH("/:tagId", core.WithBody(handler.UpdateTag))
	tagRouter.DELETE("/:tagId", handler.DeleteTag)

	// Problem-Tag association routes
	problemTagRouter := e.Group("/v1/organizations/:orgId/problems/:problemId/tags")
	problemTagRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	problemTagRouter.POST("", core.WithBody(handler.AttachTagsToProblem))
	problemTagRouter.DELETE("/:tagId", handler.DetachTagFromProblem)

	// Resource routes
	resourceRouter := e.Group("/v1/organizations/:orgId/problems/:problemId/resources")
	resourceRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	resourceRouter.POST("", core.WithBody(handler.AddResource))
	resourceRouter.GET("", handler.ListResources)
	resourceRouter.PATCH("/:resourceId", core.WithBody(handler.UpdateResource))
	resourceRouter.DELETE("/:resourceId", handler.DeleteResource)
}
