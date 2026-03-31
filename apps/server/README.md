# Coderz.space Server

[![CI](https://github.com/coderz-space/coderz.space/actions/workflows/ci.yaml/badge.svg)](https://github.com/coderz-space/coderz.space/actions/workflows/ci.yaml)

Go-based backend server for the Coderz.space bootcamp management platform.

## Quick Links

- [API Docs - Swagger UI](http://localhost:8080/swagger/index.html) (Local)
- [Docker Setup](./DOCKER.md)

## Tech Stack

- Go 1.24+
- Echo v5 (Web Framework)
- PostgreSQL 18
- SQLC (Type-safe SQL)
- JWT Authentication
- Swagger/OpenAPI

## Development

### Prerequisites

- Go 1.24+
- PostgreSQL 18
- Docker & Docker Compose (optional)
- Make

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Start PostgreSQL:
    ```bash
    make docker-up
    ```
4. Run migrations:
    ```bash
    make migrate-up
    ```
5. Generate Swagger docs:
    ```bash
    make swagger
    ```
6. Start the server:
    ```bash
    make run
    ```

### Available Commands

```bash
make build          # Build the server binary
make run            # Run development server
make swagger        # Generate Swagger documentation
make sqlc           # Generate SQLC queries
make migrate-up     # Apply database migrations
make migrate-down   # Rollback last migration
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make clear-logs     # Clear log files
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## CI/CD

The project uses GitHub Actions for continuous integration:

- **Go Server CI**: Runs tests, linting, and builds
- **Docker Build**: Validates Docker image builds

See [.github/workflows/ci.yaml](../../.github/workflows/ci.yaml) for details.

## Docker

See [DOCKER.md](./DOCKER.md) for Docker setup and deployment instructions.

Quick start with Docker Compose:

```bash
docker compose up -d
```

## Project Structure

```
apps/server/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── common/       # Shared utilities
│   ├── config/       # Configuration
│   ├── container/    # Dependency injection
│   ├── modules/      # Feature modules
│   └── routes/       # Route registration
├── db/               # Database files
│   ├── migrations/   # SQL migrations
│   └── queries/      # SQLC queries
├── swagger/          # Generated Swagger docs
└── logs/             # Application logs
```

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the server is running.

To regenerate Swagger docs after changes:

```bash
make swagger
```
