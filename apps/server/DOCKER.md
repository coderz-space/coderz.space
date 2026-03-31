# Docker Setup

## Building the Docker Image

```bash
# From the apps/server directory
docker build -t coderz-space-server:latest -f dockerfile .

# Or from the project root
docker build -t coderz-space-server:latest -f apps/server/dockerfile apps/server/
```

## Running the Container

### With docker-compose (Recommended for Development)

```bash
# Start all services (PostgreSQL + Server)
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

### Standalone Container

```bash
# Run the server container
docker run -d \
  --name coderz-server \
  -p 8080:8080 \
  -e DB_URL="postgresql://user:pass@host:5432/dbname?sslmode=disable" \
  -e JWT_SECRET="your-secret-key" \
  -e FRONTEND_ORIGIN="http://localhost:3000" \
  coderz-space-server:latest

# View logs
docker logs -f coderz-server

# Stop container
docker stop coderz-server
docker rm coderz-server
```

## Environment Variables

Required environment variables:

- `PORT` - Server port (default: 8080)
- `DB_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT token generation
- `FRONTEND_ORIGIN` - CORS allowed origin
- `ENVIRONMENT` - Environment name (development/production)

Optional:

- `LOG_LEVEL` - Logging level (default: info)
- `FILE_LOG_LEVEL` - File logging level (default: info)
- `JWT_EXPIRES` - JWT expiration time (default: 1h)

## Health Check

The container includes a health check that verifies the Swagger UI is accessible:

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' coderz-server
```

## Multi-Stage Build

The Dockerfile uses a multi-stage build:

1. **Builder stage**: Compiles the Go application
2. **Runtime stage**: Minimal Alpine image with only the binary

Benefits:

- Small image size (~20MB vs ~800MB)
- Improved security (no build tools in production)
- Faster deployment

## Production Deployment

For production, consider:

1. Using environment-specific tags
2. Implementing proper secrets management
3. Setting up health checks in your orchestrator
4. Configuring resource limits
5. Using a reverse proxy (nginx/traefik)

```bash
# Build with version tag
docker build -t coderz-space-server:v1.0.0 -f dockerfile .

# Run with resource limits
docker run -d \
  --name coderz-server \
  --memory="512m" \
  --cpus="1.0" \
  -p 8080:8080 \
  --restart unless-stopped \
  coderz-space-server:v1.0.0
```
