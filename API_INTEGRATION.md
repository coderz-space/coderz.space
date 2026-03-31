# API Integration & Docker Setup Guide

## Overview

This guide covers the API integration between the Go backend server and Next.js frontend, with Docker support for running both services independently or together.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                      Docker Network                           │
├──────────────────────┬──────────────────┬────────────────────┤
│  PostgreSQL          │  Backend API     │  Frontend Web      │
│  (Port 5432)         │  (Port 8080)     │  (Port 3000)       │
└──────────────────────┴──────────────────┴────────────────────┘
```

## Running Services

### Local Development (Without Docker)

1. **Start PostgreSQL locally** (port 5432)
   - Ensure PostgreSQL is running
   - Create database `coderz` with user `coderz-space`

2. **Start Backend Server**
   ```bash
   cd apps/server
   go mod tidy
   go run cmd/main.go
   # API runs on http://localhost:8080/api
   ```

3. **Start Frontend Web App** (in another terminal)
   ```bash
   cd apps/web
   npm install
   npm run dev
   # Web app runs on http://localhost:3000
   ```

4. **Environment Configuration**
   - Backend uses `apps/server/.env`
   - Frontend uses `apps/web/.env.local` (for local dev)
   - See `.env.example` files for reference

### Docker Deployment

#### Spin up entire stack (Web + API + Database):
```bash
docker-compose up --build
```

This will:
- Start PostgreSQL on port 5432
- Run migrations automatically
- Start API on port 8080
- Start Web on port 3000

#### Run services separately:

**Backend only:**
```bash
docker-compose up --build api postgres migrate
```

**Frontend only (requires external API):**
```bash
docker build -t coderz-web apps/web
docker run -p 3000:3000 \
  -e NEXT_PUBLIC_API_URL=http://host.docker.internal:8080/api \
  coderz-web
```

## API Integration Architecture

### Client-Side Security (Frontend)

**File:** `apps/web/services/api.ts`

Features:
- ✅ Centralized HTTP client using Axios
- ✅ Automatic auth token injection from localStorage
- ✅ Request timeouts (10 seconds)
- ✅ CORS credentials enabled
- ✅ X-Requested-With header for CSRF protection
- ✅ Automatic token refresh on 401 responses
- ✅ Custom error handling via `APIError` class
- ✅ Server-side rendering safe (no direct DOM access)

### Authentication Flow

1. **Login Request**
   - POST `/api/auth/mentee/login` with credentials
   - Backend responds with `{ token, refreshToken, mentee }`

2. **Token Storage**
   - Access token stored in `localStorage` (client-side)
   - Used automatically in `Authorization: Bearer <token>` header
   - Cleared on 401 response

3. **Protected Endpoints**
   - All subsequent requests include auth header
   - Backend middleware validates JWT
   - Invalid token triggers redirect to login

### Service Layer Integration

**File:** `apps/web/services/menteeService.ts` & `roleService.ts`

Features:
- ✅ All functions return Promises (async)
- ✅ In-memory cache with 5-minute TTL
- ✅ Automatic cache invalidation on mutations
- ✅ Graceful error handling with defaults
- ✅ TypeScript types for all responses
- ✅ Server-side rendering compatible

**Frontend Functions** → **Backend Endpoints:**

```
registerMentee()              → POST /api/auth/mentee-register
getMenteeRequests()           → GET /api/mentee-requests
updateMenteeStatus()          → PATCH /api/mentee-requests/:id
loginMentee()                 → POST /api/auth/mentee/login
loginMenteeByEmail()          → POST /api/auth/mentee/login
getMenteeQuestions()          → GET /api/mentees/:username/questions
updateQuestionProgress()      → PATCH /api/mentees/:username/questions/:questionId
updateQuestionDetails()       → PATCH /api/mentees/:username/questions/:questionId
getMenteeProfile()            → GET /api/mentees/:profileUsername/profile
getLeaderboard()              → GET /api/leaderboard
getMentorProfile()            → GET /api/mentor/profile
updateMentorProfile()         → PATCH /api/mentor/profile
selectRole()                  → POST /api/auth/select-role
getSelectedRole()             → GET /api/auth/get-role
```

## Security Best Practices Implemented

### 1. **CORS Configuration**
✅ Backend explicitly allows frontend origin only
```go
AllowOrigins:     []string{cfg.FrontendOrigin},
AllowCredentials: true,
AllowMethods:     [...specific methods...],
```

### 2. **Authentication & Authorization**
✅ JWT-based authentication
✅ Automatic token refresh handling
✅ Clear tokens on unauthorized (401) responses
✅ Tokens NOT exposed in responses headers (secure)

### 3. **Transport Security**
✅ HTTPS ready (use in production)
✅ CORS credentials enabled for secure cookies
✅ X-Requested-With header prevents CSRF
✅ Content-Type validation required

### 4. **Input Validation**
✅ Role validation in frontend service layer (defense in depth)
✅ Backend validates all inputs before database queries
✅ Error messages don't leak sensitive information

### 5. **Error Handling**
✅ Centralized error handling via `APIError` class
✅ Console warnings for debugging, not user-facing
✅ Generic error messages to prevent information leakage

### 6. **Environment Variables**
✅ Sensitive values (JWT_SECRET) never committed
✅ Different configs for local dev and Docker
✅ Production environment uses secure defaults

### 7. **Session Management**
✅ Tokens stored in localStorage (XSS-protected via CSP in production)
✅ RefreshToken for token rotation support
✅ Token expiration: 1 hour (access), 24 hours (refresh)

## Environment Variables Reference

### Frontend (`apps/web/.env.local`)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_ENVIRONMENT=development
```

### Backend (`apps/server/.env`)
```
PORT=8080
FRONTEND_ORIGIN=http://localhost:3000
JWT_SECRET=<long-secret-key>
JWT_EXPIRES=1h
```

### Docker Services
- **API connects to Database:** `postgres://coderz-space:coderz-space@postgres:5432/coderz`
- **Web connects to API:** `http://api:8080/api`
- **Frontend connects from outside:** `http://localhost:8080/api`

## Troubleshooting

### CORS Errors
**Fix:** Update `FRONTEND_ORIGIN` in backend `.env` to match frontend URL

### API Connection Failed
**Check:**
- Is backend running? `curl http://localhost:8080/api/health`
- Do hostnames match in docker-compose?
- Is firewall blocking ports?

### Docker Networking Issues
**Solution:** Services communicate via docker service names (e.g., `api`, `postgres`)
Don't use `localhost` inside Docker containers.

### Token Expiration
Clear tokens on 401, user redirected to login page automatically.

## Component Usage Example

```typescript
// components/LoginForm.tsx
import { loginMentee } from "@/services/menteeService";
import { selectRole } from "@/services/roleService";

export async function handleLogin(username: string, password: string) {
  try {
    const { token, mentee } = await loginMentee(username, password);
    // Token auto-stored by API client
    await selectRole("mentee");
    // Redirect to dashboard
  } catch (error) {
    console.error("Login failed:", error.message);
    // Show user-friendly error
  }
}
```

## Production Considerations

1. **Use HTTPS** - All API calls over HTTPS
2. **Environment Secrets** - Use secure vault for JWT_SECRET
3. **Database** - Use managed PostgreSQL service (AWS RDS, etc.)
4. **CSP Headers** - Add Content-Security-Policy headers
5. **Rate Limiting** - Implement rate limiting on backend
6. **Logging** - Monitor auth failures and errors
7. **Refresh Token Rotation** - Implement secure refresh token rotation
8. **HTTPS Enforced** - Redirect HTTP to HTTPS

## Next Steps

1. Implement remaining backend endpoints as needed
2. Add request/response logging middleware
3. Implement rate limiting
4. Add comprehensive error handling tests
5. Set up CI/CD pipeline for Docker builds
6. Configure production secrets management
