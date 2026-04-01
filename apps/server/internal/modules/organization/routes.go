package organization

import (
	"github.com/coderz-space/coderz.space/internal/common/core"
	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	orgRouter := e.Group("/v1/organizations")
	orgRouter.Use(auth.AuthMiddleware(config.JWTSecret, config.JWTExpires))

	// Organization routes
	orgRouter.POST("", core.WithBody(handler.CreateOrganization))
	orgRouter.GET("", handler.ListOrganizations)
	orgRouter.GET("/:orgId", handler.GetOrganization)
	orgRouter.PATCH("/:orgId", core.WithBody(handler.UpdateOrganization))

	// Super admin routes
	orgRouter.GET("/pending", handler.GetPendingOrganizations)
	orgRouter.POST("/:orgId/approve", handler.ApproveOrganization)

	// Member routes
	orgRouter.POST("/:orgId/members", core.WithBody(handler.AddMember))
	orgRouter.GET("/:orgId/members", handler.ListMembers)
	orgRouter.PATCH("/:orgId/members/:userId", core.WithBody(handler.UpdateMemberRole))
	orgRouter.DELETE("/:orgId/members/:userId", handler.RemoveMember)

	// Super admin cross-organization routes
	superAdminRouter := e.Group("/v1/super-admin")
	superAdminRouter.Use(auth.AuthMiddleware(config.JWTSecret, config.JWTExpires))
	superAdminRouter.GET("/organizations", handler.ListAllOrganizations)
}
