package bootcamp

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	bootcampRouter := e.Group("/v1/organizations/:orgId/bootcamps")
	bootcampRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	// Bootcamp routes
	bootcampRouter.POST("", core.WithBody(handler.CreateBootcamp))
	bootcampRouter.GET("", handler.ListBootcamps)
	bootcampRouter.GET("/:bootcampId", handler.GetBootcamp)
	bootcampRouter.PATCH("/:bootcampId", core.WithBody(handler.UpdateBootcamp))
	bootcampRouter.DELETE("/:bootcampId", handler.DeactivateBootcamp)

	// Enrollment routes
	bootcampRouter.POST("/:bootcampId/enrollments", core.WithBody(handler.EnrollMember))
	bootcampRouter.GET("/:bootcampId/enrollments", handler.ListEnrollments)
	bootcampRouter.PATCH("/:bootcampId/enrollments/:enrollmentId", core.WithBody(handler.UpdateEnrollmentRole))
	bootcampRouter.DELETE("/:bootcampId/enrollments/:enrollmentId", handler.RemoveEnrollment)
}
