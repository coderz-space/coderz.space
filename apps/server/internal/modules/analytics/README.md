# Analytics Module

## Overview

The Analytics module manages leaderboards and polls for bootcamp performance tracking and feedback collection. It provides pre-calculated leaderboard snapshots and difficulty polling functionality.

## Features

### Leaderboard Management
- **Pre-calculated Rankings**: Leaderboard entries are calculated by background jobs, not in real-time
- **Performance Metrics**: Tracks problems completed, attempted, completion rate, streak days, and score
- **Access Control**: Mentees can only view their own entry, mentors/admins can view full leaderboard
- **Pagination**: Supports offset-based pagination for large leaderboards

### Poll Management
- **Difficulty Polls**: Mentors create polls to gather feedback on problem difficulty
- **Vote Tracking**: Mentees vote on difficulty (easy, medium, hard)
- **Idempotent Voting**: PUT method allows vote creation and updates
- **Results Aggregation**: Mentors/admins can view aggregated results with percentages
- **Vote Audit**: Individual vote records available for analysis

## API Endpoints

### Leaderboard Endpoints
```
GET    /v1/bootcamps/:bootcampId/leaderboard              # Get bootcamp leaderboard
GET    /v1/bootcamps/:bootcampId/leaderboard/:enrollmentId # Get specific entry
```

### Poll Endpoints
```
POST   /v1/bootcamps/:bootcampId/polls                    # Create poll (mentor/admin)
GET    /v1/bootcamps/:bootcampId/polls                    # List polls
GET    /v1/bootcamps/:bootcampId/polls/:pollId            # Get poll details
PUT    /v1/bootcamps/:bootcampId/polls/:pollId/vote       # Vote on poll (mentee)
GET    /v1/bootcamps/:bootcampId/polls/:pollId/results    # Get results (mentor/admin)
GET    /v1/bootcamps/:bootcampId/polls/:pollId/votes      # Get individual votes (mentor/admin)
```

## Authorization Rules

### Leaderboard
- All enrolled users can view the full leaderboard
- Mentees can only view their own detailed entry
- Mentors/admins can view all entries

### Polls
- Mentors/admins can create polls
- All enrolled users can view polls
- Only mentees can vote on polls
- Only mentors/admins/super_admins can view results and individual votes

## Data Models

### Leaderboard Entry
```go
type LeaderboardEntryData struct {
    ID                   UUID
    BootcampID           UUID
    BootcampEnrollmentID UUID
    Rank                 int32
    ProblemsCompleted    int32
    ProblemsAttempted    int32
    CompletionRate       string
    StreakDays           int32
    Score                int32
    CalculatedAt         string
    Name                 string
    AvatarURL            string
}
```

### Poll
```go
type PollData struct {
    ID           UUID
    BootcampID   UUID
    ProblemID    UUID
    Question     string
    CreatedBy    UUID
    CreatedAt    string
    ProblemTitle string
    MyVote       string  // User's vote if they voted
}
```

### Vote
```go
type VoteData struct {
    ID        UUID
    PollID    UUID
    VoterID   UUID  // Bootcamp enrollment ID
    Vote      string  // easy, medium, hard
    CreatedAt string
}
```

## Validation Rules

### Poll Creation
- Question: 10-240 characters
- Problem ID: Valid UUID, must exist
- User must be mentor/admin

### Poll Voting
- Vote: Must be one of: easy, medium, hard
- User must be mentee
- User must be enrolled in poll's bootcamp

## Implementation Notes

### Leaderboard Calculation
- Leaderboard entries are pre-calculated by background jobs
- The API serves snapshot data, not real-time calculations
- Use `UpsertLeaderboardEntry` service method for background job updates
- Entries include `calculated_at` timestamp for freshness tracking

### Poll Voting
- Uses PUT method for idempotent vote creation/update
- Returns 201 for first vote, 200 for updates
- Voter ID is the bootcamp enrollment ID, not user ID
- Prevents exposure of internal user identifiers

### Access Control
- All operations require bootcamp enrollment verification
- Role-based filtering enforced at service layer
- Mentees restricted from viewing aggregated results
- CSRF protection required for cookie-based authentication

## Dependencies

- **SQLC Queries**: `analytics.sql`
- **Auth Middleware**: JWT token validation
- **Common Utils**: UUID parsing, timestamp formatting
- **Validator**: Struct validation

## Testing

See `analytics_test.go` for comprehensive test coverage including:
- Leaderboard retrieval and pagination
- Poll creation and listing
- Vote casting and updates
- Results aggregation
- Access control enforcement
- Multi-tenant isolation

## Future Enhancements

- Real-time leaderboard updates via WebSocket
- Advanced analytics (time-to-completion, difficulty trends)
- Poll templates and reusable questions
- Export functionality for results
- Leaderboard filtering by time period
