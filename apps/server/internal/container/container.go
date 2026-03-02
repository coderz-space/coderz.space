package container

import (
	"github.com/DSAwithGautam/CodeConquerers/internal/config"
	"go.uber.org/zap"
)

type Container struct {

	// configuration
	Config *config.Config
	Logger *zap.Logger

	// services

	// handlers

	// repositories

	// database
}

func NewContainer(config *config.Config, logger *zap.Logger) (*Container, error) {
	container := &Container{
		Config: config,
		Logger: logger,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	return nil
}
