package container

import (
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/db"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Container struct {

	// configuration
	Config *config.Config
	Logger *zap.Logger

	// auth
	AuthHandler *auth.Handler
	AuthService *auth.Service

	// DB
	DB *pgxpool.Pool
}

func NewContainer(config *config.Config, logger *zap.Logger) (*Container, error) {

	authService := auth.NewService()
	authHandler := auth.NewHandler(authService)

	db, err := db.InitDB(config)
	if err != nil {
		return nil, err
	}

	container := &Container{
		Config:      config,
		Logger:      logger,
		AuthHandler: authHandler,
		AuthService: authService,
		DB:          db,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	c.DB.Close()
	return nil
}
