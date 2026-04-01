package auth

import (
	"testing"

	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
)

func TestNewService(t *testing.T) {
	queries := &db.Queries{}
	cfg := &config.Config{}

	service := NewService(queries, cfg)

	if service == nil {
		t.Fatal("Expected NewService to return a non-nil Service instance")
	}

	if service.queries != queries {
		t.Errorf("Expected queries to be %p, got %p", queries, service.queries)
	}

	if service.config != cfg {
		t.Errorf("Expected config to be %p, got %p", cfg, service.config)
	}
}
