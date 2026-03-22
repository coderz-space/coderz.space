package main

import (
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/logger"
	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware"
	config "github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/container"
	"github.com/DSAwithGautam/Coderz.space/internal/routes"
	_ "github.com/DSAwithGautam/Coderz.space/swagger" // Import generated docs
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)



// @title Coderz.space API
// @version 1.0
// @description This is a server for Coderz.space
// @host localhost:8080
// @BasePath /api
func main() {

	cfg := config.LoadConfig()
	logger.Initialize(cfg)
	defer logger.Sync()

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

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// register routes
	router := e.Group("/api")
	routes.RegisterRoutes(router, di)

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
