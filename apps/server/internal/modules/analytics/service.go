package analytics

import (
	"context"
	"errors"
	"fmt"

	"github.com/coderz-space/coderz.space/internal/common/utils"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		queries: db.New(pool),
		pool:    pool,
	}
}

// Leaderboard Service Methods

// GetBootcampLeaderboard retrieves pre-calculated leaderboard entries for a bootcamp
func (s *Service) GetBootcampLeaderboard(ctx context.Context, bootcampID pgtype.UUID, page, limit int) ([]LeaderboardEntryData, int, error) {
	// Fetch leaderboard entries (pre-calculated snapshots)
	entries, err := s.queries.GetLeaderboardByBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	total := len(entries)
	offset := (page - 1) * limit
	end := offset + limit

	if offset >= total {
		return []LeaderboardEntryData{}, total, nil
	}

	if end > total {
		end = total
	}

	// Map to response data
	data := make([]LeaderboardEntryData, 0, end-offset)
	for i := offset; i < end; i++ {
		data = append(data, mapLeaderboardEntryToData(&entries[i]))
	}

	return data, total, nil
}

// GetLeaderboardEntry retrieves a single leaderboard entry with access control
func (s *Service) GetLeaderboardEntry(ctx context.Context, bootcampID, enrollmentID pgtype.UUID, userRole string, memberID pgtype.UUID) (*LeaderboardEntryResponse, error) {
	// Fetch all entries to find the specific one
	entries, err := s.queries.GetLeaderboardByBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, err
	}

	// Find the entry
	var entry *db.GetLeaderboardByBootcampRow
	for i := range entries {
		if entries[i].BootcampEnrollmentID == enrollmentID {
			entry = &entries[i]
			break
		}
	}

	if entry == nil {
		return nil, errors.New("ENTRY_NOT_FOUND")
	}

	// Access control: mentees can only view their own entry
	if userRole == "mentee" {
		// Get enrollment for this member to verify ownership
		memberEnrollment, err := s.queries.GetEnrollmentByMemberID(ctx, db.GetEnrollmentByMemberIDParams{
			OrganizationMemberID: memberID,
			BootcampID:           bootcampID,
		})
		if err != nil {
			return nil, errors.New("ACCESS_DENIED")
		}

		// Check if the requested enrollment belongs to this member
		if memberEnrollment.ID != enrollmentID {
			return nil, errors.New("ACCESS_DENIED")
		}
	}

	return &LeaderboardEntryResponse{
		Success: true,
		Data:    mapLeaderboardEntryToData(entry),
	}, nil
}

// UpsertLeaderboardEntry creates or updates a leaderboard entry (for background jobs)
func (s *Service) UpsertLeaderboardEntry(ctx context.Context, bootcampID pgtype.UUID, req UpsertLeaderboardEntryRequest) (*LeaderboardEntryData, error) {
	enrollmentID, err := utils.StringToUUID(req.BootcampEnrollmentID)
	if err != nil {
		return nil, errors.New("INVALID_ENROLLMENT_ID")
	}

	entry, err := s.queries.UpsertLeaderboardEntry(ctx, db.UpsertLeaderboardEntryParams{
		BootcampID:           bootcampID,
		BootcampEnrollmentID: enrollmentID,
		ProblemsCompleted:    req.ProblemsCompleted,
		ProblemsAttempted:    req.ProblemsAttempted,
		CompletionRate:       float32(req.CompletionRate),
		StreakDays:           req.StreakDays,
		Score:                req.Score,
		Rank:                 req.Rank,
	})
	if err != nil {
		return nil, err
	}

	// Map to response (without user details since this is for background jobs)
	return &LeaderboardEntryData{
		ID:                   entry.ID,
		BootcampID:           entry.BootcampID,
		BootcampEnrollmentID: entry.BootcampEnrollmentID,
		Rank:                 entry.Rank,
		ProblemsCompleted:    entry.ProblemsCompleted,
		ProblemsAttempted:    entry.ProblemsAttempted,
		CompletionRate:       fmt.Sprintf("%.2f", entry.CompletionRate),
		StreakDays:           entry.StreakDays,
		Score:                entry.Score,
		CalculatedAt:         utils.FormatTimestamp(entry.CalculatedAt),
	}, nil
}

// Poll Service Methods

// CreatePoll creates a new poll for a problem in a bootcamp
func (s *Service) CreatePoll(ctx context.Context, bootcampID pgtype.UUID, req CreatePollRequest, createdBy pgtype.UUID) (*PollResponse, error) {
	problemID, err := utils.StringToUUID(req.ProblemID)
	if err != nil {
		return nil, errors.New("INVALID_PROBLEM_ID")
	}

	// Validate problem exists
	// The database foreign key constraint will handle validation
	// If problem doesn't exist, the insert will fail

	poll, err := s.queries.CreatePoll(ctx, db.CreatePollParams{
		BootcampID: bootcampID,
		ProblemID:  problemID,
		Question:   req.Question,
		CreatedBy:  createdBy,
	})
	if err != nil {
		// Check if it's a foreign key violation (problem not found)
		if err.Error() == "ERROR: insert or update on table \"polls\" violates foreign key constraint (SQLSTATE 23503)" {
			return nil, errors.New("PROBLEM_NOT_FOUND")
		}
		return nil, err
	}

	return &PollResponse{
		Success: true,
		Data: PollData{
			ID:         poll.ID,
			BootcampID: poll.BootcampID,
			ProblemID:  poll.ProblemID,
			Question:   poll.Question,
			CreatedBy:  poll.CreatedBy,
			CreatedAt:  utils.FormatTimestamp(poll.CreatedAt),
		},
	}, nil
}

// ListPolls retrieves polls for a bootcamp with optional problem filtering
func (s *Service) ListPolls(ctx context.Context, bootcampID pgtype.UUID, problemIDStr string, voterID pgtype.UUID, page, limit int) ([]PollData, int, error) {
	// Fetch polls
	polls, err := s.queries.ListPollsByBootcamp(ctx, bootcampID)
	if err != nil {
		return nil, 0, err
	}

	// Filter by problem ID if provided
	var filtered []db.ListPollsByBootcampRow
	if problemIDStr != "" {
		problemID, err := utils.StringToUUID(problemIDStr)
		if err == nil {
			for i := range polls {
				if polls[i].ProblemID == problemID {
					filtered = append(filtered, polls[i])
				}
			}
			polls = filtered
		}
	}

	// Calculate pagination
	total := len(polls)
	offset := (page - 1) * limit
	end := offset + limit

	if offset >= total {
		return []PollData{}, total, nil
	}

	if end > total {
		end = total
	}

	// Map to response data with user's vote
	data := make([]PollData, 0, end-offset)
	for i := offset; i < end; i++ {
		pollData := mapPollToData(&polls[i])

		// Get user's vote for this poll if they voted
		vote, err := s.queries.GetUserVoteForPoll(ctx, db.GetUserVoteForPollParams{
			PollID:  polls[i].ID,
			VoterID: voterID,
		})
		if err == nil {
			pollData.MyVote = string(vote.Vote)
		}

		data = append(data, pollData)
	}

	return data, total, nil
}

// GetPoll retrieves a single poll with user's vote state
func (s *Service) GetPoll(ctx context.Context, bootcampID, pollID, voterID pgtype.UUID) (*PollResponse, error) {
	poll, err := s.queries.GetPoll(ctx, pollID)
	if err != nil {
		return nil, errors.New("POLL_NOT_FOUND")
	}

	// Validate poll belongs to bootcamp
	if poll.BootcampID != bootcampID {
		return nil, errors.New("POLL_NOT_FOUND")
	}

	pollData := PollData{
		ID:         poll.ID,
		BootcampID: poll.BootcampID,
		ProblemID:  poll.ProblemID,
		Question:   poll.Question,
		CreatedBy:  poll.CreatedBy,
		CreatedAt:  utils.FormatTimestamp(poll.CreatedAt),
	}

	// Get user's vote for this poll if they voted
	vote, err := s.queries.GetUserVoteForPoll(ctx, db.GetUserVoteForPollParams{
		PollID:  pollID,
		VoterID: voterID,
	})
	if err == nil {
		pollData.MyVote = string(vote.Vote)
	}

	return &PollResponse{
		Success: true,
		Data:    pollData,
	}, nil
}

// VotePoll casts or updates a vote on a poll
func (s *Service) VotePoll(ctx context.Context, pollID, voterID pgtype.UUID, vote string) (*VoteResponse, bool, error) {
	// Validate poll exists
	_, err := s.queries.GetPoll(ctx, pollID)
	if err != nil {
		return nil, false, errors.New("POLL_NOT_FOUND")
	}

	// Check if vote already exists (for determining status code)
	voteExists, err := s.queries.CheckVoteExists(ctx, db.CheckVoteExistsParams{
		PollID:  pollID,
		VoterID: voterID,
	})
	if err != nil {
		return nil, false, err
	}

	isNew := !voteExists

	// Cast vote (upsert)
	voteRecord, err := s.queries.CastPollVote(ctx, db.CastPollVoteParams{
		PollID:  pollID,
		VoterID: voterID,
		Vote:    db.PollVoteValue(vote),
	})
	if err != nil {
		return nil, false, err
	}

	return &VoteResponse{
		Success: true,
		Data: VoteData{
			ID:        voteRecord.ID,
			PollID:    voteRecord.PollID,
			VoterID:   voteRecord.VoterID,
			Vote:      string(voteRecord.Vote),
			CreatedAt: utils.FormatTimestamp(voteRecord.CreatedAt),
		},
	}, isNew, nil
}

// GetPollResults retrieves aggregated poll results
func (s *Service) GetPollResults(ctx context.Context, pollID pgtype.UUID) (*PollResultsResponse, error) {
	// Validate poll exists
	_, err := s.queries.GetPoll(ctx, pollID)
	if err != nil {
		return nil, errors.New("POLL_NOT_FOUND")
	}

	// Get vote counts
	results, err := s.queries.GetPollResults(ctx, pollID)
	if err != nil {
		return nil, err
	}

	// Aggregate results
	var totalVotes int32
	voteBreakdown := make(map[string]int32)
	percentBreakup := make(map[string]float64)

	for _, result := range results {
		count := int32(result.VoteCount) // #nosec G115 - VoteCount is from database count
		voteBreakdown[string(result.Vote)] = count
		totalVotes += count
	}

	// Calculate percentages
	if totalVotes > 0 {
		for vote, count := range voteBreakdown {
			percentBreakup[vote] = float64(count) / float64(totalVotes) * 100
		}
	}

	return &PollResultsResponse{
		Success: true,
		Data: PollResultsData{
			TotalVotes:     totalVotes,
			EasyCount:      voteBreakdown["easy"],
			MediumCount:    voteBreakdown["medium"],
			HardCount:      voteBreakdown["hard"],
			EasyPercent:    percentBreakup["easy"],
			MediumPercent:  percentBreakup["medium"],
			HardPercent:    percentBreakup["hard"],
			VoteBreakdown:  voteBreakdown,
			PercentBreakup: percentBreakup,
		},
	}, nil
}

// GetPollVotes retrieves individual vote records with optional filtering
func (s *Service) GetPollVotes(ctx context.Context, pollID pgtype.UUID, voteFilter string, page, limit int) ([]VoteData, int, error) {
	// Validate poll exists
	_, err := s.queries.GetPoll(ctx, pollID)
	if err != nil {
		return nil, 0, errors.New("POLL_NOT_FOUND")
	}

	// Count total votes
	total, err := s.queries.CountPollVotesByPoll(ctx, db.CountPollVotesByPollParams{
		PollID:  pollID,
		Column2: voteFilter,
	})
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch votes with pagination
	votes, err := s.queries.ListPollVotesByPoll(ctx, db.ListPollVotesByPollParams{
		PollID:  pollID,
		Column2: voteFilter,
		Limit:   int32(limit),  // #nosec G115 - limit is bounded by max 100
		Offset:  int32(offset), // #nosec G115 - offset is calculated from bounded values
	})
	if err != nil {
		return nil, 0, err
	}

	// Map to response data
	data := make([]VoteData, len(votes))
	for i := range votes {
		data[i] = VoteData{
			ID:        votes[i].ID,
			PollID:    votes[i].PollID,
			VoterID:   votes[i].VoterID,
			Vote:      string(votes[i].Vote),
			CreatedAt: utils.FormatTimestamp(votes[i].CreatedAt),
		}
	}

	return data, int(total), nil // #nosec G115 - total is from database count
}

// Helper Methods

// GetMemberIDByUserAndBootcamp retrieves the organization member ID for a user in a bootcamp
func (s *Service) GetMemberIDByUserAndBootcamp(ctx context.Context, userID, bootcampID pgtype.UUID) (pgtype.UUID, error) {
	memberID, err := s.queries.GetMemberIDByUserID(ctx, db.GetMemberIDByUserIDParams{
		UserID:     userID,
		BootcampID: bootcampID,
	})
	if err != nil {
		return pgtype.UUID{}, errors.New("MEMBER_NOT_FOUND")
	}
	return memberID, nil
}

// GetEnrollmentIDByUserAndBootcamp retrieves the bootcamp enrollment ID for a user
func (s *Service) GetEnrollmentIDByUserAndBootcamp(ctx context.Context, userID, bootcampID pgtype.UUID) (pgtype.UUID, error) {
	enrollmentID, err := s.queries.GetEnrollmentIDByUserID(ctx, db.GetEnrollmentIDByUserIDParams{
		UserID:     userID,
		BootcampID: bootcampID,
	})
	if err != nil {
		return pgtype.UUID{}, errors.New("ENROLLMENT_NOT_FOUND")
	}
	return enrollmentID, nil
}

// Mapping Functions

func mapLeaderboardEntryToData(entry *db.GetLeaderboardByBootcampRow) LeaderboardEntryData {
	return LeaderboardEntryData{
		ID:                   entry.ID,
		BootcampID:           entry.BootcampID,
		BootcampEnrollmentID: entry.BootcampEnrollmentID,
		Rank:                 entry.Rank,
		ProblemsCompleted:    entry.ProblemsCompleted,
		ProblemsAttempted:    entry.ProblemsAttempted,
		CompletionRate:       fmt.Sprintf("%.2f", entry.CompletionRate),
		StreakDays:           entry.StreakDays,
		Score:                entry.Score,
		CalculatedAt:         utils.FormatTimestamp(entry.CalculatedAt),
		Name:                 entry.Name,
		AvatarURL:            formatNullableText(entry.AvatarUrl),
	}
}

func mapPollToData(poll *db.ListPollsByBootcampRow) PollData {
	return PollData{
		ID:           poll.ID,
		BootcampID:   poll.BootcampID,
		ProblemID:    poll.ProblemID,
		Question:     poll.Question,
		CreatedBy:    poll.CreatedBy,
		CreatedAt:    utils.FormatTimestamp(poll.CreatedAt),
		ProblemTitle: poll.ProblemTitle,
	}
}

func formatNullableText(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}
