package auth

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	authRouter := e.Group("/v1/auth")
	authRouter.POST("/login", core.WithBody(handler.Login))
	authRouter.POST("/signup", core.WithBody(handler.Signup))
	authRouter.POST("/refresh", handler.Refresh)
}

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	authRouter := e.Group("/v1/auth")
	authRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

	authRouter.GET("/me", handler.Me)
	authRouter.POST("/logout", handler.Logout)
}
