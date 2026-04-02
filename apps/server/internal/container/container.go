package container

import (
	"github.com/coderz-space/coderz.space/internal/config"
	"github.com/coderz-space/coderz.space/internal/db"
	db_sqlc "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/coderz-space/coderz.space/internal/common/email"
	"github.com/coderz-space/coderz.space/internal/modules/analytics"
	"github.com/coderz-space/coderz.space/internal/modules/app"
	"github.com/coderz-space/coderz.space/internal/modules/assignment"
	"github.com/coderz-space/coderz.space/internal/modules/auth"
	"github.com/coderz-space/coderz.space/internal/modules/bootcamp"
	"github.com/coderz-space/coderz.space/internal/modules/organization"
	"github.com/coderz-space/coderz.space/internal/modules/problem"
	"github.com/coderz-space/coderz.space/internal/modules/progress"
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

	// app facade
	AppHandler *app.Handler
	AppService *app.Service

	// DB
	DB *pgxpool.Pool
}

func NewContainer(config *config.Config, logger *zap.Logger) (*Container, error) {

	pool, err := db.InitDB(config)
	if err != nil {
		return nil, err
	}

	queries := db_sqlc.New(pool)

	// Initialize email service
	emailService := email.NewService(config)

	// Initialize auth module
	authService := auth.NewService(queries, config, emailService)
	authHandler := auth.NewHandler(authService)

	// Initialize organization module
	organizationService := organization.NewService(queries, config, pool)
	organizationHandler := organization.NewHandler(organizationService)

	// Initialize bootcamp module
	bootcampService := bootcamp.NewService(queries, config, pool)
	bootcampHandler := bootcamp.NewHandler(bootcampService)

	// Initialize problem module
	problemService := problem.NewService(queries, config, pool)
	problemHandler := problem.NewHandler(problemService)

	// Initialize assignment module
	assignmentService := assignment.NewService(pool, queries)
	assignmentHandler := assignment.NewHandler(assignmentService)

	// Initialize progress module
	progressService := progress.NewService(pool)
	progressHandler := progress.NewHandler(progressService)

	// Initialize analytics module
	analyticsService := analytics.NewService(pool)
	analyticsHandler := analytics.NewHandler(analyticsService)

	// Initialize app facade module
	appService := app.NewService(pool)
	appHandler := app.NewHandler(appService)

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
		AppHandler:          appHandler,
		AppService:          appService,
		DB:                  pool,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	c.DB.Close()
	return nil
}
