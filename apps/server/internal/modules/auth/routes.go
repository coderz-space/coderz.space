package auth

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/core"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/labstack/echo/v5"
)

func RegisterPublicRoutes(e *echo.Group, handler *Handler) {
	authRouter := e.Group("/v1/auth")
	authRouter.POST("/signin", core.WithBody(handler.SignIn))
	authRouter.POST("/register", core.WithBody(handler.Register))

}

func RegisterProtectedRoutes(e *echo.Group, handler *Handler, config *config.Config) {
	authRouter := e.Group("/v1/auth")
	authRouter.Use(auth.AuthMiddleware(config.JWT_SECRET, config.JWT_EXPIRES))

}
