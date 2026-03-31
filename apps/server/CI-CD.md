# CI/CD Pipeline Documentation

## Overview

The Coderz.space server uses GitHub Actions for continuous integration and Docker for containerization.

## CI Pipeline

### Workflow: `.github/workflows/ci.yaml`

The CI pipeline runs on:

- All pull requests
- Pushes to `main`, `master`, `dev`, and `prod` branches

### Jobs

#### 1. Go Server CI (`go-server`)

**Services:**

- PostgreSQL 18 (for integration tests)

**Steps:**

1. Checkout code
2. Setup Go 1.25.x with dependency caching
3. Install and verify dependencies
4. Run `go vet` for static analysis
5. Run `staticcheck` for additional linting
6. Run `golangci-lint` with comprehensive checks
7. Setup test environment with database
8. Run database migrations
9. Execute tests with race detection and coverage
10. Upload coverage to Codecov
11. Build server binary
12. Generate and verify Swagger documentation

**Linting Tools:**

- `go vet`: Built-in Go static analyzer
- `staticcheck`: Advanced static analysis
- `golangci-lint`: Meta-linter running multiple linters

**Configuration:**

- Linter config: `.golangci.yml`
- Timeout: 5 minutes
- Coverage: Atomic mode with race detection

#### 2. Docker Build (`docker-build`)

**Dependencies:** Requires `go-server` job to pass

**Steps:**

1. Checkout code
2. Setup Docker Buildx
3. Build Docker image with caching
4. Validate image builds successfully

**Optimizations:**

- GitHub Actions cache for layers
- Multi-stage build for minimal image size

## Docker Setup

### Dockerfile

**Location:** `apps/server/dockerfile`

**Build Strategy:** Multi-stage build

- **Stage 1 (builder):** Go 1.25-alpine with build tools
- **Stage 2 (runtime):** Minimal Alpine with only the binary

**Features:**

- Non-root user execution
- Health check endpoint
- Swagger documentation included
- ~20MB final image size

### Docker Compose

**Location:** `apps/server/docker-compose.yml`

**Services:**

1. **postgres**: PostgreSQL 18 database
2. **migrate**: Database migration runner
3. **server**: Go application server

**Features:**

- Health checks for all services
- Automatic migration on startup
- Environment variable configuration
- Volume persistence for database

## Local Development

### Running with Docker Compose

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f server

# Stop services
docker compose down
```

### Running without Docker

```bash
# Start PostgreSQL only
make docker-up

# Run migrations
make migrate-up

# Start server
make run
```

## Environment Variables

### Required

- `PORT`: Server port (default: 8080)
- `DB_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing secret
- `FRONTEND_ORIGIN`: CORS allowed origin

### Optional

- `ENVIRONMENT`: Environment name (development/production)
- `LOG_LEVEL`: Application log level (info/debug/warn/error)
- `FILE_LOG_LEVEL`: File log level
- `JWT_EXPIRES`: Token expiration time (default: 1h)
- `MAX_DB_CONNS`: Max database connections (default: 10)

See `.env.example` for complete list.

## Testing Strategy

### Unit Tests

```bash
go test ./...
```

### Integration Tests

```bash
# Requires PostgreSQL running
make docker-up
make migrate-up
go test -v ./...
```

### Coverage

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Deployment

### Building for Production

```bash
# Build Docker image
docker build -t coderz-space-server:v1.0.0 -f dockerfile .

# Push to registry
docker tag coderz-space-server:v1.0.0 registry.example.com/coderz-space-server:v1.0.0
docker push registry.example.com/coderz-space-server:v1.0.0
```

### Production Considerations

1. **Secrets Management**
    - Use environment-specific secrets
    - Never commit `.env` files
    - Use secret management tools (Vault, AWS Secrets Manager)

2. **Database Migrations**
    - Run migrations before deploying new version
    - Test migrations on staging first
    - Keep rollback scripts ready

3. **Health Checks**
    - Endpoint: `/swagger/index.html`
    - Interval: 30s
    - Timeout: 3s
    - Start period: 10s

4. **Resource Limits**
    - Memory: 512MB recommended
    - CPU: 1.0 core recommended
    - Adjust based on load

5. **Monitoring**
    - Application logs in `logs/` directory
    - Structured JSON logging
    - Log rotation with lumberjack

## Troubleshooting

### CI Failures

**Linting errors:**

```bash
# Run locally
golangci-lint run --timeout=5m
```

**Test failures:**

```bash
# Run with verbose output
go test -v ./...
```

**Build failures:**

```bash
# Verify dependencies
go mod verify
go mod tidy
```

### Docker Issues

**Build failures:**

```bash
# Check Dockerfile syntax
docker build -f dockerfile .
```

**Container won't start:**

```bash
# Check logs
docker logs coderz-space-server

# Check health
docker inspect --format='{{.State.Health.Status}}' coderz-space-server
```

**Database connection issues:**

```bash
# Verify PostgreSQL is running
docker compose ps

# Check database logs
docker compose logs postgres
```

## Maintenance

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update specific dependency
go get -u github.com/labstack/echo/v5@latest
go mod tidy
```

### Regenerating Swagger Docs

```bash
make swagger
# or
swag init -o ./swagger --parseDependency --parseInternal -g cmd/main.go
```

### Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir db/migrations -seq migration_name

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down
```
