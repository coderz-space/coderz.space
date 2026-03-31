package container

import (
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/db"
	db_sqlc "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/organization"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/progress"
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

	// organization
	OrganizationHandler *organization.Handler
	OrganizationService *organization.Service

	// progress (doubts)
	ProgressHandler *progress.Handler
	ProgressService *progress.Service

	// DB
	DB *pgxpool.Pool
}

func NewContainer(config *config.Config, logger *zap.Logger) (*Container, error) {

	pool, err := db.InitDB(config)
	if err != nil {
		return nil, err
	}

	queries := db_sqlc.New(pool)

	// Initialize auth module
	authService := auth.NewService(queries, config)
	authHandler := auth.NewHandler(authService)

	// Initialize organization module
	organizationService := organization.NewService(queries, config, pool)
	organizationHandler := organization.NewHandler(organizationService)

	// Initialize progress module
	progressService := progress.NewService(pool)
	progressHandler := progress.NewHandler(progressService)

	container := &Container{
		Config:              config,
		Logger:              logger,
		AuthHandler:         authHandler,
		AuthService:         authService,
		OrganizationHandler: organizationHandler,
		OrganizationService: organizationService,
		ProgressHandler:     progressHandler,
		ProgressService:     progressService,
		DB:                  pool,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	c.DB.Close()
	return nil
}
