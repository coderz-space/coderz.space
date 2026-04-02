package main

import (
	"net/http"
	"time"

	"github.com/coderz-space/coderz.space/internal/common/logger"
	"github.com/coderz-space/coderz.space/internal/common/middleware"
	"github.com/coderz-space/coderz.space/internal/common/middleware/timeout"
	config "github.com/coderz-space/coderz.space/internal/config"
	"github.com/coderz-space/coderz.space/internal/container"
	"github.com/coderz-space/coderz.space/internal/routes"
	_ "github.com/coderz-space/coderz.space/swagger" // Import generated docs
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

// @title Coderz.space Bootcamp Management API
// @version 1.0
// @description Comprehensive bootcamp management platform API with multi-tenant architecture and role-based access control
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@coderz.space

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Organizations
// @tag.description Organization management endpoints

// @tag.name Organization Members
// @tag.description Organization member management endpoints

// @tag.name Bootcamps
// @tag.description Bootcamp lifecycle management endpoints

// @tag.name Bootcamp Enrollments
// @tag.description Bootcamp enrollment management endpoints
func main() {

	cfg := config.LoadConfig()
	logger.Initialize(cfg)
	defer func() {
		_ = logger.Sync() // Best effort sync on shutdown
	}()

	di, err := container.NewContainer(cfg, logger.Logger)
	if err != nil {
		logger.Fatal("failed to create container", zap.Error(err))
	}
	defer di.Close()

	e := echo.New()

	// middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Request-Id"},
		AllowCredentials: true,
	}))
	e.Use(middleware.ZapLogger())
	e.Use(middleware.Recovery())
	e.Use(timeout.TimeoutMiddleware(30 * time.Second)) // 30 second timeout to prevent resource exhaustion

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// register routes
	router := e.Group("/api")
	routes.RegisterRoutes(router, di)

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
