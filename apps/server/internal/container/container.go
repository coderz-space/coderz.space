package container

import (
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/DSAwithGautam/Coderz.space/internal/db"
	db_sqlc "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/mentorship"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/organization"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/profile"
	"github.com/DSAwithGautam/Coderz.space/internal/modules/tasks"
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

	// mentorship
	MentorshipHandler *mentorship.Handler
	MentorshipService *mentorship.Service

	// tasks
	TasksHandler *tasks.Handler
	TasksService *tasks.Service

	// profile
	ProfileHandler *profile.Handler
	ProfileService *profile.Service

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

	// Initialize mentorship module
	mentorshipService := mentorship.NewService(queries)
	mentorshipHandler := mentorship.NewHandler(mentorshipService)

	// Initialize tasks module
	tasksService := tasks.NewService(queries)
	tasksHandler := tasks.NewHandler(tasksService)

	// Initialize profile module
	profileService := profile.NewService(queries)
	profileHandler := profile.NewHandler(profileService)

	container := &Container{
		Config:              config,
		Logger:              logger,
		AuthHandler:         authHandler,
		AuthService:         authService,
		OrganizationHandler: organizationHandler,
		OrganizationService: organizationService,
		MentorshipHandler:   mentorshipHandler,
		MentorshipService:   mentorshipService,
		TasksHandler:        tasksHandler,
		TasksService:        tasksService,
		ProfileHandler:      profileHandler,
		ProfileService:      profileService,
		DB:                  pool,
	}
	return container, nil
}

// close all the connections and cleans up all resources
func (c *Container) Close() error {
	c.DB.Close()
	return nil
}
