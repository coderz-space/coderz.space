package db

import (
	"context"
	"fmt"

	"github.com/DSAwithGautam/Coderz.space/internal/common/logger"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	// configuration
	config, err := pgxpool.ParseConfig(cfg.DB_URL)
	if err != nil {
		return nil, err
	}

	// set pool configuration
	if cfg.MaxDBConns > 0 {
		config.MaxConns = int32(cfg.MaxDBConns)
	}
	if cfg.MinDBConns > 0 {
		config.MinConns = int32(cfg.MinDBConns)
	}
	config.MaxConnLifetime = cfg.MaxDBConnLifetime
	config.MaxConnIdleTime = cfg.MaxDBConnIdleTime

	// create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	logger.Info("database connection established")
	return pool, nil
}
