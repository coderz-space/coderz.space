package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAllowsMissingDotEnvWhenEnvironmentIsProvided(t *testing.T) {
	t.Setenv("APP_NAME", "Coderz_Space")
	t.Setenv("VERSION", "0.1.0")
	t.Setenv("PORT", "8080")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_EXPIRES", "1h")
	t.Setenv("FRONTEND_ORIGIN", "http://localhost:3000")
	t.Setenv("DB_URL", "postgres://localhost:5432/coderz?sslmode=disable")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("REFRESH_TOKEN_EXPIRES", "24h")
	t.Setenv("MAX_DB_CONN_LIFETIME", "1h")
	t.Setenv("MAX_DB_CONN_IDLE_TIME", "30m")
	t.Setenv("MAX_DB_CONNS", "10")
	t.Setenv("MIN_DB_CONNS", "2")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("FILE_LOG_LEVEL", "info")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(wd)))
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(root)
	})

	cfg := LoadConfig()
	if cfg.AppName != "Coderz_Space" {
		t.Fatalf("unexpected app name: %s", cfg.AppName)
	}
}
