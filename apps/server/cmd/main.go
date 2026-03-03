package main

import (
	"net/http"
	"time"

	"github.com/DSAwithGautam/CodeConquerers/internal/common/logger"
	"github.com/DSAwithGautam/CodeConquerers/internal/common/middleware"
	config "github.com/DSAwithGautam/CodeConquerers/internal/config"
	"github.com/DSAwithGautam/CodeConquerers/internal/container"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

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

	// health check api :
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status":      "ok",
			"service":     cfg.AppName,
			"version":     cfg.Version,
			"environment": cfg.Environment,
			"database":    "nil",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
