# Docker Debugging Guide

## Environment Variables Inside Docker

When services run in Docker, they can communicate via service names:

```yaml
# docker-compose.yml defines:
services:
  api:      # Service name = hostname
  web:      # Can access api as: http://api:8080
  postgres: # Can access db as: postgres:5432
```

## Environment Variable Mapping

### Frontend Service Names
Inside Docker container, frontend connects to:
```
NEXT_PUBLIC_API_URL=http://api:8080/api
```

From your laptop browser, connect to:
```
http://localhost:3000 → calls → http://localhost:8080/api
```

### Backend Service Names
Inside Docker container, backend connects to:
```
DB_URL=postgres://user:pass@postgres:5432/coderz
```

From your laptop psql client, connect to:
```
psql -h localhost -p 5432 -U coderz-space coderz
```

## Verification Commands

### Check if containers are running
```bash
docker ps
```

Expected output:
```
coderz-api
coderz-web
coderz-postgres
```

### Check container logs
```bash
# API logs
docker-compose logs api

# Web logs
docker-compose logs web

# Database logs
docker-compose logs postgres

# Follow logs in real-time
docker-compose logs -f api
```

### Test API from container
```bash
# From your laptop
curl http://localhost:8080/api/health

# Expected response
{"status":"ok","timestamp":"2026-03-31T..."}
```

### Test network connectivity inside containers
```bash
# Open shell in API container
docker-compose exec api sh

# Inside container, test DB connection
nc -zv postgres 5432  # Should show: postgres:5432 open

# Test API health
wget http://localhost:8080/api/health -O -
```

## Common Issues

### Issue: "Connection refused" to API from frontend

**Cause:** Frontend using `http://localhost:8080` instead of `http://api:8080` inside Docker

**Solution:**
Check `.env.production`:
```
# Wrong for Docker
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Correct for Docker
NEXT_PUBLIC_API_URL=http://api:8080/api
```

**Rebuild:** `docker-compose up --build web`

### Issue: Database migrations not running

**Check migration logs:**
```bash
docker-compose logs migrate
```

**Common causes:**
- PostgreSQL not healthy yet (wait for health check)
- Wrong DB connection string
- Missing migration files

**Fix:**
```bash
docker-compose down -v  # Remove volume
docker-compose up --build  # Rebuild everything
```

### Issue: Port already in use

**Cause:** Another service using port 3000, 8080, or 5432

**Solution:**
```bash
# Find what's using the port (Linux/Mac)
lsof -i :8080

# Kill it
kill <PID>

# Or change port in docker-compose.yml
# ports:
#   - "8081:8080"  # Changed from 8080
```

### Issue: Containers keep restarting

**Check logs:**
```bash
docker-compose logs <service-name>
```

**Common causes:**
- Database not initialized
- Wrong environment variables
- Port conflicts
- Out of memory

**Debug:**
```bash
# Run container in foreground to see errors
docker-compose run --rm api sh

# Inside container, run server manually
./server  # See actual error
```

### Issue: Frontend can't see API even though it's running

**Check:**
1. Is API health check passing?
   ```bash
   docker-compose ps
   # Look for "healthy" status
   ```

2. Is web connected to network?
   ```bash
   docker network inspect coderz-network
   # Should list both 'api' and 'web' containers
   ```

3. Can web reach API from container?
   ```bash
   docker-compose exec web wget -O - http://api:8080/api/health
   ```

**Solution:**
```bash
docker-compose down
docker-compose up --build
```

## Environment Variable Debugging

### Print environment inside container
```bash
# In API container
docker-compose exec api env | grep -E "API|DB|FRONTEND"

# In web container
docker-compose exec web env | grep -E "NEXT_PUBLIC"
```

### Verify environment variables loaded
Check container startup logs:
```bash
docker-compose logs api | grep -E "PORT|ORIGIN|DATABASE"
```

### Override environment at runtime
```bash
docker run -e NEXT_PUBLIC_API_URL=http://example.com coderz-web
```

## Performance Debugging

### Container resource usage
```bash
docker stats  # See CPU, memory, network usage

# Monitor specific container
docker stats coderz-api
```

### Slow startup?
```bash
# Check when each step completed
docker-compose logs --timestamps api

# Timings:
# 1. Build image (~30s)
# 2. Start database (~5s)
# 3. Run migrations (~5s)
# 4. Start API (~2s)
# 5. Start web (~15s)
```

## Volume & Persistence

### Check volume status
```bash
docker volume ls | grep coderz
docker volume inspect coderz-postgres-data
```

### Remove volume (WARNING: deletes data!)
```bash
docker-compose down -v
```

### Backup database from Docker
```bash
docker-compose exec postgres pg_dump -U coderz-space coderz > backup.sql
```

### Restore database
```bash
cat backup.sql | docker-compose exec -T postgres psql -U coderz-space coderz
```

## Network Debugging

### Inspect docker network
```bash
docker network inspect coderz-network
```

Shows all containers connected and their IP addresses.

### Test DNS resolution inside container
```bash
docker-compose exec api nslookup postgres
# Should resolve to 172.x.x.x
```

### Check exposed ports
```bash
docker ps --format "table {{.Names}}\t{{.Ports}}"
```

## Security Verification

### Check CORS headers
```bash
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     http://localhost:8080/api/health -v
```

Should see `Access-Control-Allow-Origin: http://localhost:3000`

### Verify JWT validation
1. Login to get token
2. Test with wrong token
3. Should get 401 Unauthorized

### Check auth flow
```bash
# Login
TOKEN=$(curl -X POST http://localhost:8080/api/auth/mentee/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}' | jq -r '.token')

# Use token
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/mentees/testuser/profile
```

## Rebuild & Restart

### Rebuild everything
```bash
docker-compose down
docker-compose up --build
```

### Rebuild specific service
```bash
docker-compose up --build api  # Rebuild only API
```

### Hard reset (remove everything)
```bash
docker-compose down -v          # Stop & remove volumes
docker system prune -a          # Clean unused images
docker-compose up --build       # Fresh start
```

## Production Debugging

### Enable debug mode
Add to `.env`:
```
LOG_LEVEL=debug
```

Rebuild:
```bash
docker-compose up --build
```

### View request/response in logs
API logs should show:
- Request method & path
- Response status code
- Processing time

Frontend logs (browser console):
- API call details
- Response data or errors

### Monitor API metrics
```bash
# Check response times
docker-compose logs api | grep "duration"

# Find slow requests (>1s)
docker-compose logs api | grep "duration.*[1-9][0-9][0-9][0-9]ms"
```

## Extracting Logs for Support

```bash
# Save all logs to file
docker-compose logs > debug.log

# Just API logs
docker-compose logs api > api.log

# With timestamps
docker-compose logs --timestamps > debug_time.log

# Follow in real-time
docker-compose logs -f
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `docker-compose up` | Start all services |
| `docker-compose down` | Stop all services |
| `docker-compose ps` | List running containers |
| `docker-compose logs api` | View API logs |
| `docker-compose exec api sh` | Shell into API container |
| `docker-compose build` | Rebuild images |
| `docker stats` | Monitor resource usage |
| `docker system prune -a` | Clean up everything |

## Getting Help

1. Check logs first: `docker-compose logs -f`
2. Review [QUICKSTART.md](./QUICKSTART.md) troubleshooting
3. Check [API_INTEGRATION.md](./API_INTEGRATION.md) for architecture
4. Verify all containers healthy: `docker-compose ps`
5. Try full reset: `docker-compose down -v && docker-compose up --build`
