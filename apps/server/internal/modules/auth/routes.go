package auth

import (
	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	authRouter := e.Group("/v1/auth")
	authRouter.POST("/login", handler.Login)
	authRouter.POST("/signup", handler.Signup)
	authRouter.POST("/refresh", handler.Refresh)
	authRouter.POST("/forgot-password", handler.ForgotPassword)
	authRouter.POST("/reset-password", handler.ResetPassword)
}

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	authRouter := e.Group("/v1/auth")
	authRouter.Use(auth.AuthMiddleware(config.JWTSecret, config.JWTExpires))

	authRouter.GET("/me", handler.Me)
	authRouter.POST("/logout", handler.Logout)
}
