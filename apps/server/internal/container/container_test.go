package container

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestContainer_Close(t *testing.T) {
	// Create a dummy pool
	config, err := pgxpool.ParseConfig("postgres://localhost:5432/postgres")
	assert.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	assert.NoError(t, err)

	// Create a container with the dummy pool
	c := &Container{
		DB: pool,
	}

	// Call Close and assert no error
	err = c.Close()
	assert.NoError(t, err)

	// Additional verification: we can't easily check if pool is closed directly
	// without calling a method that requires it to be open, or relying on internals.
	// We can try to Ping it. It should fail since the pool is closed.
	err = pool.Ping(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "closed pool")
}
