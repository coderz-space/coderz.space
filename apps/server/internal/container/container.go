package container

import (
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/db"
	db_sqlc "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/analytics"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/assignment"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/bootcamp"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/organization"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/problem"
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

	// bootcamp
	BootcampHandler *bootcamp.Handler
	BootcampService *bootcamp.Service

	// problem
	ProblemHandler *problem.Handler
	ProblemService *problem.Service

	// assignment
	AssignmentHandler *assignment.Handler
	AssignmentService *assignment.Service

	// progress (doubts)
	ProgressHandler *progress.Handler
	ProgressService *progress.Service

	// analytics
	AnalyticsHandler *analytics.Handler
	AnalyticsService *analytics.Service

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

	// Initialize bootcamp module
	bootcampService := bootcamp.NewService(pool)
	bootcampHandler := bootcamp.NewHandler(bootcampService)

	// Initialize problem module
	problemService := problem.NewService(pool)
	problemHandler := problem.NewHandler(problemService)

	// Initialize assignment module
	assignmentService := assignment.NewService(pool)
	assignmentHandler := assignment.NewHandler(assignmentService)

	// Initialize progress module
	progressService := progress.NewService(pool)
	progressHandler := progress.NewHandler(progressService)

	// Initialize analytics module
	analyticsService := analytics.NewService(pool)
	analyticsHandler := analytics.NewHandler(analyticsService)

	container := &Container{
		Config:              config,
		Logger:              logger,
		AuthHandler:         authHandler,
		AuthService:         authService,
		OrganizationHandler: organizationHandler,
		OrganizationService: organizationService,
		BootcampHandler:     bootcampHandler,
		BootcampService:     bootcampService,
		ProblemHandler:      problemHandler,
		ProblemService:      problemService,
		AssignmentHandler:   assignmentHandler,
		AssignmentService:   assignmentService,
		ProgressHandler:     progressHandler,
		ProgressService:     progressService,
		AnalyticsHandler:    analyticsHandler,
		AnalyticsService:    analyticsService,
		DB:                  pool,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	c.DB.Close()
	return nil
}
