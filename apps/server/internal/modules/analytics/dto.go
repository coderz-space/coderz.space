package analytics

import "github.com/jackc/pgx/v5/pgtype"

// Leaderboard DTOs

// LeaderboardEntryData represents a single leaderboard entry
// @Description Leaderboard entry with user details and performance metrics
type LeaderboardEntryData struct {
	ID                   pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BootcampID           pgtype.UUID `json:"bootcampId" example:"660e8400-e29b-41d4-a716-446655440000"`
	BootcampEnrollmentID pgtype.UUID `json:"bootcampEnrollmentId" example:"770e8400-e29b-41d4-a716-446655440000"`
	Rank                 int32       `json:"rank" example:"1"`
	ProblemsCompleted    int32       `json:"problemsCompleted" example:"25"`
	ProblemsAttempted    int32       `json:"problemsAttempted" example:"30"`
	CompletionRate       string      `json:"completionRate" example:"83.33"`
	StreakDays           int32       `json:"streakDays" example:"7"`
	Score                int32       `json:"score" example:"850"`
	CalculatedAt         string      `json:"calculatedAt" example:"2024-01-15T10:30:00Z"`
	Name                 string      `json:"name" example:"John Doe"`
	AvatarURL            string      `json:"avatarUrl,omitempty" example:"https://example.com/avatar.jpg"`
}

// LeaderboardResponse represents a list of leaderboard entries
// @Description Response containing leaderboard entries with pagination
type LeaderboardResponse struct {
	Data    []LeaderboardEntryData `json:"data"`
	Meta    *OffsetPagination      `json:"meta,omitempty"`
	Success bool                   `json:"success" example:"true"`
}

// LeaderboardEntryResponse represents a single leaderboard entry response
// @Description Response containing a single leaderboard entry
type LeaderboardEntryResponse struct {
	Data    LeaderboardEntryData `json:"data"`
	Success bool                 `json:"success" example:"true"`
}

// UpsertLeaderboardEntryRequest represents the request to upsert a leaderboard entry
// @Description Request body for upserting a leaderboard entry (background job use)
type UpsertLeaderboardEntryRequest struct {
	BootcampEnrollmentID string  `json:"bootcampEnrollmentId" validate:"required,uuid" example:"770e8400-e29b-41d4-a716-446655440000"`
	ProblemsCompleted    int32   `json:"problemsCompleted" validate:"required,min=0" example:"25"`
	ProblemsAttempted    int32   `json:"problemsAttempted" validate:"required,min=0" example:"30"`
	CompletionRate       float64 `json:"completionRate" validate:"required,min=0,max=100" example:"83.33"`
	StreakDays           int32   `json:"streakDays" validate:"required,min=0" example:"7"`
	Score                int32   `json:"score" validate:"required,min=0" example:"850"`
	Rank                 int32   `json:"rank" validate:"required,min=1" example:"1"`
}

// Poll DTOs

// CreatePollRequest represents the request body for creating a poll
// @Description Request body for creating a poll on a problem
type CreatePollRequest struct {
	ProblemID string `json:"problemId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Question  string `json:"question" validate:"required,min=10,max=240" example:"How difficult did you find this problem?"`
}

// VotePollRequest represents the request body for voting on a poll
// @Description Request body for casting or updating a vote on a poll
type VotePollRequest struct {
	Vote string `json:"vote" validate:"required,oneof=easy medium hard" example:"medium"`
}

// PollData represents poll details
// @Description Poll details with problem information
type PollData struct {
	ID           pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BootcampID   pgtype.UUID `json:"bootcampId" example:"660e8400-e29b-41d4-a716-446655440000"`
	ProblemID    pgtype.UUID `json:"problemId" example:"770e8400-e29b-41d4-a716-446655440000"`
	Question     string      `json:"question" example:"How difficult did you find this problem?"`
	CreatedBy    pgtype.UUID `json:"createdBy" example:"880e8400-e29b-41d4-a716-446655440000"`
	CreatedAt    string      `json:"createdAt" example:"2024-01-15T09:00:00Z"`
	ProblemTitle string      `json:"problemTitle,omitempty" example:"Two Sum"`
	MyVote       string      `json:"myVote,omitempty" example:"medium"`
}

// PollResponse represents a single poll response
// @Description Response containing a single poll
type PollResponse struct {
	Data    PollData `json:"data"`
	Success bool     `json:"success" example:"true"`
}

// PollListResponse represents a list of polls
// @Description Response containing a list of polls with pagination
type PollListResponse struct {
	Data    []PollData        `json:"data"`
	Meta    *OffsetPagination `json:"meta,omitempty"`
	Success bool              `json:"success" example:"true"`
}

// VoteData represents a poll vote
// @Description Poll vote details
type VoteData struct {
	ID        pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PollID    pgtype.UUID `json:"pollId" example:"660e8400-e29b-41d4-a716-446655440000"`
	VoterID   pgtype.UUID `json:"voterId" example:"770e8400-e29b-41d4-a716-446655440000"`
	Vote      string      `json:"vote" example:"medium"`
	CreatedAt string      `json:"createdAt" example:"2024-01-15T09:00:00Z"`
}

// VoteResponse represents a single vote response
// @Description Response containing a single vote
type VoteResponse struct {
	Data    VoteData `json:"data"`
	Success bool     `json:"success" example:"true"`
}

// PollResultsData represents aggregated poll results
// @Description Aggregated poll results with vote counts and percentages
type PollResultsData struct {
	TotalVotes     int32              `json:"totalVotes" example:"100"`
	EasyCount      int32              `json:"easyCount" example:"20"`
	MediumCount    int32              `json:"mediumCount" example:"50"`
	HardCount      int32              `json:"hardCount" example:"30"`
	EasyPercent    float64            `json:"easyPercent" example:"20.0"`
	MediumPercent  float64            `json:"mediumPercent" example:"50.0"`
	HardPercent    float64            `json:"hardPercent" example:"30.0"`
	VoteBreakdown  map[string]int32   `json:"voteBreakdown"`
	PercentBreakup map[string]float64 `json:"percentBreakup"`
}

// PollResultsResponse represents poll results response
// @Description Response containing aggregated poll results
type PollResultsResponse struct {
	Data    PollResultsData `json:"data"`
	Success bool            `json:"success" example:"true"`
}

// PollVotesResponse represents a list of individual votes
// @Description Response containing individual poll votes with pagination
type PollVotesResponse struct {
	Data    []VoteData        `json:"data"`
	Meta    *OffsetPagination `json:"meta,omitempty"`
	Success bool              `json:"success" example:"true"`
}

// OffsetPagination represents offset-based pagination metadata
// @Description Offset-based pagination metadata
type OffsetPagination struct {
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
	Total int `json:"total" example:"100"`
}

// GenericResponse represents a generic success response
// @Description Generic success response
type GenericResponse struct {
	Data    map[string]any `json:"data"`
	Success bool           `json:"success" example:"true"`
}
