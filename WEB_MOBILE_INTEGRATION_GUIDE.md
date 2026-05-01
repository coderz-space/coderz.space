# WEB & MOBILE INTEGRATION QUICK START

**For**: Frontend developers (Web & Mobile teams)  
**Updated**: April 2026  
**Server Status**: 95% Ready (3-5 days to full deployment)

---

## 🔑 THE 3 THINGS YOU NEED TO KNOW RIGHT NOW

### 1️⃣ Server is Ready (Almost)

✅ **Status**: All 48/50 endpoints working  
⚠️ **TODO**: Email service for password reset (non-blocking)  
**What this means**: You can start building against the API TODAY. Password reset will be fixed this week.

### 2️⃣ How to Get Running

```bash
# API runs on: http://localhost:8080/api
# Swagger docs at: http://localhost:8080/swagger/index.html
# Start with: make docker-up && make migrate-up && make run
```

### 3️⃣ Authentication Format

```
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

---

## 🎯 REQUIRED ENDPOINTS FOR WEB (Next.js)

### Authentication Flow

```
POST /v1/auth/signup           (Public) → Register user
POST /v1/auth/login            (Public) → Authenticate & get tokens
POST /v1/auth/refresh          (Public) → Refresh access token
GET  /v1/auth/me               (Protected) → Get current user profile
POST /v1/auth/logout           (Protected) → Logout
POST /v1/auth/password-reset   (Public) → Reset password (⚠️ email TODO)
```

### Dashboard & Navigation

```
GET  /v1/app/context                              → Get user's org/bootcamp context
GET  /v1/app/bootcamps                            → List user's bootcamps
GET  /v1/organizations                            → List user's organizations
GET  /v1/organizations/{orgId}                    → Get org details
GET  /v1/organizations/{orgId}/bootcamps          → List org's bootcamps
GET  /v1/organizations/{orgId}/bootcamps/{bootcampId} → Get bootcamp details
```

### Mentor Dashboard

```
GET  /v1/organizations/{orgId}/problems                           → List problems
GET  /v1/organizations/{orgId}/problems/{problemId}               → Get problem details
POST /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups      → Create assignment
GET  /v1/organizations/{orgId}/bootcamps/{bootcampId}/assignment-groups      → List assignments
GET  /v1/bootcamps/{bootcampId}/enrollments                       → List bootcamp enrollments
PATCH /v1/bootcamps/{bootcampId}/enrollments/{enrollmentId}      → Update enrollment
```

### Mentee Dashboard

```
GET  /v1/bootcamps/{bootcampId}/leaderboard                      → View leaderboard
GET  /v1/bootcamps/{bootcampId}/leaderboard/{enrollmentId}       → View own rank
GET  /v1/bootcamps/{bootcampId}/assignments-for-mentee           → Get assigned problems
PUT  /v1/assignments/{assignmentId}/problems/{problemId}         → Submit solution
GET  /v1/bootcamps/{bootcampId}/polls                            → View polls
PUT  /v1/bootcamps/{bootcampId}/polls/{pollId}/vote              → Vote on poll
```

### Questions/Doubts

```
POST /v1/doubts                           → Raise doubt
GET  /v1/doubts                           → List doubts (Mentor)
GET  /v1/doubts/me                        → Get my doubts (Mentee)
GET  /v1/doubts/{doubtId}                 → Get doubt details
PATCH /v1/doubts/{doubtId}/resolve        → Resolve doubt (Mentor)
```

---

## 📱 REQUIRED ENDPOINTS FOR MOBILE (React Native)

### Authentication (Same as Web)

```
POST /v1/auth/signup
POST /v1/auth/login
POST /v1/auth/refresh
GET  /v1/auth/me
POST /v1/auth/logout
```

### Core Navigation

```
GET  /v1/app/context                              → Dashboard init
GET  /v1/bootcamps/{bootcampId}/assignments-for-mentee  → Daily assignments
```

### Problem Solving

```
GET  /v1/bootcamps/{bootcampId}/assignments/{assignmentId}/problems
GET  /v1/assignments/{assignmentId}/problems/{problemId}
PUT  /v1/assignments/{assignmentId}/problems/{problemId}    → Submit/update status
POST /v1/assignments/{assignmentId}/problems/{problemId}/skip → Skip problem
```

### Doubt Management

```
POST /v1/doubts                      → Raise doubt on problem
GET  /v1/doubts/me                   → Get my doubts
GET  /v1/doubts/{doubtId}            → Get doubt details
PATCH /v1/doubts/{doubtId}/resolve   → See mentor's answer
```

### Progress Tracking

```
GET  /v1/bootcamps/{bootcampId}/leaderboard  → View leaderboard
GET  /v1/bootcamps/{bootcampId}/polls        → Polls for feedback
PUT  /v1/bootcamps/{bootcampId}/polls/{pollId}/vote
```

---

## 🔐 AUTHENTICATION FLOW

### Step-by-Step for Web/Mobile

**1. User Signup**

```javascript
POST /v1/auth/signup
Body: {
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure_password_123"
}
Response: {
  "success": true,
  "data": {
    "id": "uuid",
    "email": "john@example.com",
    "name": "John Doe",
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "new_refresh_token"
  }
}
Headers: Set-Cookie: refresh_token=eyJ...; HttpOnly; SameSite=Strict
```

**2. Store Tokens**

- Access Token: In-memory or sessionStorage (short-lived, usually 1h)
- Refresh Token: HttpOnly cookie (automatic, browser handles)

**3. Make API Requests**

```javascript
// Include access token
fetch("/api/v1/auth/me", {
  headers: {
    Authorization: "Bearer " + accessToken,
  },
});
```

**4. Handle Token Expiration**
When you get 401:

```javascript
POST / v1 / auth / refresh;
// Browser automatically includes refresh_token cookie
// Response gives new accessToken
```

**5. Logout**

```javascript
POST / v1 / auth / logout;
// Clears server-side refresh token
// Frontend should clear accessToken
```

---

## 🚨 ERROR HANDLING CHECKLIST

### Status Codes You'll See

| Code | Meaning      | Handle By                                       |
| ---- | ------------ | ----------------------------------------------- |
| 200  | Success      | Use data                                        |
| 201  | Created      | Use data, update UI                             |
| 400  | Bad request  | Show user validation error                      |
| 401  | Unauthorized | Refresh token or redirect to login              |
| 403  | Forbidden    | Show "You don't have permission"                |
| 404  | Not found    | Show "Resource not found"                       |
| 409  | Conflict     | Show specific error (slug exists, email exists) |
| 429  | Rate limited | Wait 60 seconds, retry (doubts endpoint)        |
| 500  | Server error | Retry later, show error to user                 |

### Response Format

```json
{
  "success": false,
  "status": "BAD_REQUEST",
  "message": "VALIDATION_ERROR",
  "data": null,
  "meta": null
}
```

### Common Error Codes

```
INVALID_REQUEST_BODY
VALIDATION_ERROR
VALIDATION_FAILED
INVALID_TOKEN_CLAIMS
INVALID_USER_ID
UNAUTHORIZED
NOT_ORGANIZATION_MEMBER
NOT_ENROLLED_IN_BOOTCAMP
NOT_FOUND
CONFLICT
SLUG_ALREADY_EXISTS
EMAIL_ALREADY_EXISTS
ORGANIZATION_NOT_APPROVED
BOOTCAMP_NOT_FOUND
ASSIGNMENT_PROBLEM_NOT_FOUND
INVALID_ASSIGNMENT_PROBLEM_ID
```

---

## 📊 PAGINATION PATTERNS

### Offset-based (Most endpoints)

```
GET /v1/organizations?page=1&limit=20

Response.meta: {
  "page": 1,
  "limit": 20,
  "total": 150
}

// Client calculates total pages = Math.ceil(total / limit)
```

### Cursor-based (Doubts endpoint)

```
GET /v1/doubts?cursor=&limit=20&bootcampId=xxx

Response includes:
"meta": {
  "nextCursor": "cursor_value" or null if no more
}

// nextCursor = null means you've reached the end
```

---

## 🌱 SAMPLE API CALL FLOW FOR WEB

### Mentee Dashboard Flow

```javascript
// 1. Get user context
GET /api/v1/app/context
→ Returns list of bootcamps user is enrolled in

// 2. Get bootcamp details
GET /api/v1/organizations/{orgId}/bootcamps/{bootcampId}
→ Get start date, end date, etc.

// 3. Get assignments for today
GET /api/v1/bootcamps/{bootcampId}/assignments-for-mentee
→ Get problems assigned for mentee

// 4. Get leaderboard
GET /api/v1/bootcamps/{bootcampId}/leaderboard?page=1&limit=20
→ Show user their rank

// 5. Get doubts
GET /api/v1/doubts/me?bootcampId={id}&resolved=false
→ Show unresolved doubts

// 6. Get polls
GET /api/v1/bootcamps/{bootcampId}/polls
→ Show difficulty polls for feedback
```

---

## 🌱 SAMPLE API CALL FLOW FOR MOBILE

### Problem-Solving Flow

```javascript
// 1. Get assignments
GET /api/v1/bootcamps/{bootcampId}/assignments-for-mentee

// 2. Get problem details
GET /api/v1/assignments/{assignmentId}/problems/{problemId}
→ {
  "title": "Two Sum",
  "description": "Find two numbers that add up...",
  "difficulty": "easy",
  "resources": ["link1", "link2"],
  "status": "pending"
}

// 3. User submits solution
PUT /api/v1/assignments/{assignmentId}/problems/{problemId}
Body: { "status": "completed", "code": "..." }
→ Update problem status

// 4. Raise doubt if stuck
POST /api/v1/doubts
Body: { "assignmentProblemId": "...", "message": "..." }

// 5. View mentor's response
GET /api/v1/doubts/{doubtId}
→ See resolution_note from mentor
```

---

## ✅ READY-TO-USE FEATURES

### Feature Completeness Matrix

| Feature          | Web | Mobile      | Status |
| ---------------- | --- | ----------- | ------ |
| Authentication   | ✅  | ✅          | Ready  |
| Organizations    | ✅  | ⚠️ Optional | Ready  |
| Bootcamps        | ✅  | ✅          | Ready  |
| Problems         | ✅  | ✅          | Ready  |
| Assignments      | ✅  | ✅          | Ready  |
| Leaderboard      | ✅  | ✅          | Ready  |
| Doubts/Questions | ✅  | ✅          | Ready  |
| Polls            | ✅  | ✅          | Ready  |
| Admin Dashboard  | ✅  | ⚠️ Optional | Ready  |

---

## 🎓 TESTING BEFORE LAUNCH

### What Frontend Devs Should Test

1. **Authentication**
   - [ ] Signup with valid email
   - [ ] Signup with taken email (should fail)
   - [ ] Login with correct password
   - [ ] Login with wrong password (should fail)
   - [ ] Refresh token after expiration
   - [ ] Access protected endpoint without token (should get 401)

2. **Data Fetching**
   - [ ] Get bootcamp list (verify pagination)
   - [ ] Get problems with filters
   - [ ] Get leaderboard with sorting
   - [ ] Get assignments for mentee

3. **Creating Data**
   - [ ] Submit problem solution
   - [ ] Raise doubt/question
   - [ ] Vote on poll

4. **Error Cases**
   - [ ] Network offline
   - [ ] Server returns 500
   - [ ] Token expires mid-request
   - [ ] Submit invalid data (validation error)

---

## 📋 DEPENDENCY CHECKLIST FOR FRONTEND TEAMS

### Web (Next.js)

- [ ] `axios` or `fetch` for API calls
- [ ] JWT token management library (js-jwt)
- [ ] Form validation library
- [ ] Error handling middleware
- [ ] State management (Redux, Zustand, etc.)
- [ ] Environment variables (.env)

### Mobile (React Native)

- [ ] `axios` or `fetch` for API calls
- [ ] Secure token storage (react-native-keychain)
- [ ] Form validation
- [ ] Error handling
- [ ] State management (Redux, Zustand, Recoil)
- [ ] Environment configuration

---

## 🔗 IMPORTANT LINKS

| Resource              | Link                                              |
| --------------------- | ------------------------------------------------- |
| Swagger UI            | http://localhost:8080/swagger/index.html          |
| Health Check          | http://localhost:8080/api/health                  |
| GitHub Repo           | [Your repo URL]                                   |
| API Base              | http://localhost:8080/api                         |
| Main README           | /apps/server/README.md                            |
| Progress Module Docs  | /apps/server/internal/modules/progress/README.md  |
| Analytics Module Docs | /apps/server/internal/modules/analytics/README.md |

---

## 🚨 KNOWN LIMITATIONS & WORKAROUNDS

### Email Service (TODO this week)

**Issue**: Password reset emails not sending yet  
**Workaround**: Manual password reset for now  
**ETA Fix**: 2-3 days

### Rate Limiting on Doubts

**Limit**: 10 requests per minute per user  
**Workaround**: Show users remaining attempts in UI  
**Example**: `POST /v1/doubts` can be called max 10x per 60 seconds

All other features are production-ready with no known limitations.

---

## 💡 BEST PRACTICES

### 1. Always Include Authorization Header

```javascript
headers: {
  'Authorization': 'Bearer ' + tokens.accessToken,
  'Content-Type': 'application/json'
}
```

### 2. Handle Pagination Correctly

```javascript
// Don't assume all data comes at once
// Implement pagination UI/infinite scroll
for page=1; page <= totalPages; page++ {
  fetch(`/v1/endpoint?page=${page}&limit=20`)
}
```

### 3. Cache When Possible

```javascript
// Problems, bootcamps don't change frequently
// Cache with 1-hour TTL or user-triggered refresh
```

### 4. Show Loading & Error States

```javascript
// API responses take time
// Show spinner while loading
// Show error message if request fails
```

### 5. Validate Input Before Sending

```javascript
// Validate email format, password strength, etc. on client
// Reduces unnecessary API calls
```

---

## ✨ NEXT STEPS FOR YOUR TEAM

### Week 1: Setup & Planning

- [ ] Access Swagger UI and explore endpoints
- [ ] Run server locally: `make docker-up && make run`
- [ ] Create API client wrapper/service layer
- [ ] Plan authentication flow
- [ ] Design state management structure

### Week 2: Build Core Features

- [ ] Implement authentication (signup/login/logout)
- [ ] Build dashboard/navigation
- [ ] Fetch and display data (bootcamps, assignments)
- [ ] Setup error handling

### Week 3: Implement Features

- [ ] Build problem-solving interface
- [ ] Implement progress tracking
- [ ] Build doubts/comments system
- [ ] Add leaderboard visualization

### Week 4: Polish & Test

- [ ] Full end-to-end testing
- [ ] Performance optimization
- [ ] UI/UX refinement
- [ ] Security review

---

## 📞 QUICK REFERENCE CHEATSHEET

```bash
# Start server
cd apps/server && make docker-up && make migrate-up && make run

# Health check
curl http://localhost:8080/api/health

# Signup & login
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"pass123"}'

# Get user profile (with token)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get bootcamps
curl "http://localhost:8080/api/v1/app/bootcamps" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

**Questions?** Check the auto-generated docs:

- API_COMPLETENESS_ANALYSIS.md
- INTEGRATION_PLAN.md
- Module READMEs (progress/, analytics/)

**Last Updated**: April 2026
