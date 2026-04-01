package container

import (
	"testing"

	"github.com/coderz-space/coderz.space/internal/config"
	"go.uber.org/zap"
)

func TestNewContainerError(t *testing.T) {
	cfg := &config.Config{
		DBURL: "invalid_postgres_url://host:port/db",
	}
	logger := zap.NewNop()

	c, err := NewContainer(cfg, logger)

	if err == nil {
		if c != nil {
			c.Close()
		}
		t.Fatalf("expected error when initializing container with invalid DB URL, got nil")
	}

	if c != nil {
		t.Fatalf("expected container to be nil on error, got %v", c)
	}
}
