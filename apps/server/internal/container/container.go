package container

import (
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"go.uber.org/zap"
)

type Container struct {

	// configuration
	Config *config.Config
	Logger *zap.Logger

	// auth
	AuthHandler *auth.Handler
	AuthService *auth.Service
}

func NewContainer(config *config.Config, logger *zap.Logger) (*Container, error) {

	authService := auth.NewService()
	authHandler := auth.NewHandler(authService, config)

	container := &Container{
		Config:      config,
		Logger:      logger,
		AuthHandler: authHandler,
		AuthService: authService,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	return nil
}
