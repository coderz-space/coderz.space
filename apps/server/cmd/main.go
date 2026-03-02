package main

import (
	"net/http"
	"time"

	"github.com/DSAwithGautam/CodeConquerers/internal/common/logger"
	config "github.com/DSAwithGautam/CodeConquerers/internal/config"
	"github.com/DSAwithGautam/CodeConquerers/internal/container"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.NewLogger(cfg)
	defer logger.Sync()

	di, err := container.NewContainer(cfg, logger)
	if err != nil {
		logger.Fatal("failed to create container", zap.Error(err))
	}
	defer di.Close()

	e := echo.New()

	// health check :
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
