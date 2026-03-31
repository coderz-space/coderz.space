# Quick Start - API Integration

## 📋 What's Been Set Up

- ✅ Secure HTTP client (`services/api.ts`)
- ✅ API-integrated services (`roleService.ts`, `menteeService.ts`)
- ✅ JWT authentication with auto-token injection
- ✅ Environment configuration for local & Docker
- ✅ Docker Compose with Web + API + PostgreSQL
- ✅ Security best practices implemented
- ✅ Comprehensive documentation

## 🚀 Get Started

### Option 1: Full Stack in Docker (Recommended)

```bash
# From project root
docker-compose up --build

# Services will be available at:
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080/api
# Health:   http://localhost:8080/api/health
```

### Option 2: Local Development

**Terminal 1 - Backend:**
```bash
cd apps/server
go run cmd/main.go
# Runs on http://localhost:8080/api
```

**Terminal 2 - Frontend:**
```bash
cd apps/web
npm install
npm run dev
# Runs on http://localhost:3000
```

**Start PostgreSQL independently:**
- Docker: `docker run -p 5432:5432 -e POSTGRES_USER=coderz-space -e POSTGRES_PASSWORD=coderz-space postgres:18`
- Or use local PostgreSQL installation

## 📝 Environment Setup

### For Local Development Edit

**`apps/web/.env.local`:**
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_ENVIRONMENT=development
```

**`apps/server/.env`:**
- Already configured for localhost
- Update `FRONTEND_ORIGIN` if using different frontend URL

## 🧪 Test the Integration

### Check API Health
```bash
curl http://localhost:8080/api/health
# Response: {"status":"ok","timestamp":"2026-03-31T..."}
```

### Test Frontend Connection
1. Open http://localhost:3000
2. Browser console should not show CORS errors
3. Try logging in - requests should go to backend

## 📚 Documentation

- **API Integration Guide:** [API_INTEGRATION.md](./API_INTEGRATION.md)
- **Security Audit:** [SECURITY_AUDIT.md](./SECURITY_AUDIT.md)

## 🔧 Docker Commands

```bash
# Start everything
docker-compose up --build

# Stop everything
docker-compose down

# View logs
docker-compose logs -f api
docker-compose logs -f web

# Rebuild specific service
docker-compose up --build api

# Run backend only
docker-compose up postgres migrate api

# Clean up everything (including data)
docker-compose down -v
```

## 🔌 API Endpoints (Implemented)

All endpoints below are integrated in frontend services:

### Authentication
- `POST /api/auth/mentee-register` - Register new mentee
- `POST /api/auth/mentee/login` - Login mentee
- `POST /api/auth/select-role` - Select user role
- `GET /api/auth/get-role` - Get selected role

### Mentee Management
- `GET /api/mentee-requests` - Get all mentee requests (admin)
- `PATCH /api/mentee-requests/:id` - Update mentee status
- `DELETE /api/mentee-requests/:id` - Delete mentee
- `GET /api/mentees/:username/profile` - Get mentee profile
- `PATCH /api/mentees/:username/profile` - Update mentee profile
- `PATCH /api/mentees/:username/password` - Change password

### Questions & Progress
- `GET /api/mentees/:username/questions` - Get questions
- `PATCH /api/mentees/:username/questions/:questionId` - Update progress/notes
- `GET /api/mentees/:username/questions/:questionId` - Get question detail

### Leaderboard
- `GET /api/leaderboard` - Get mentee rankings

### Mentor
- `GET /api/mentor/profile` - Get mentor profile
- `PATCH /api/mentor/profile` - Update mentor profile
- `PATCH /api/mentor/password` - Change password

### Health
- `GET /api/health` - Health check

## 🛡️ Security Features

✅ **CORS Protection** - Only frontend can access API
✅ **JWT Authentication** - Secure token-based auth
✅ **Auto Token Injection** - No manual header management
✅ **Centralized Error Handling** - Generic error messages
✅ **Cache Layer** - Reduced API load with TTL
✅ **Type Safety** - Full TypeScript support

## 🚨 Common Issues

| Issue | Solution |
|-------|----------|
| CORS Error | Check FRONTEND_ORIGIN in `.env` |
| Cannot connect to DB | Ensure PostgreSQL is running |
| Port already in use | `docker-compose down` or change ports |
| API not responding | Check logs: `docker-compose logs api` |

## 📦 Dependencies Added

- **Frontend:** `axios@^1.7.0` (HTTP client)
- **Backend:** Already complete

Install frontend dependencies:
```bash
cd apps/web
npm install
```

## ✅ Features Keeping Existing UI

All frontend components remain unchanged:
- UI components, layouts, and styling intact
- Only service layer implementations updated
- Backward compatible with existing component code
- No breaking changes to component APIs

## 🎯 Next: Implement Backend Endpoints

The frontend is now ready. Backend should implement the API endpoints mapped in `API_INTEGRATION.md`.

Start with these core endpoints:
1. Auth endpoints (login, register, role selection)
2. Mentee questions endpoint
3. Profile endpoints
4. Leaderboard endpoint

## 📞 Support

See `API_INTEGRATION.md` for detailed troubleshooting and architecture diagrams.
