package analytics

import (
	"net/http"

	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/coderz-space/coderz.space/internal/common/validator"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Leaderboard Handlers

// GetBootcampLeaderboard godoc
// @Summary Get bootcamp leaderboard
// @Description Retrieve pre-calculated leaderboard rankings for a bootcamp. Returns snapshot data without real-time recalculation. User must be enrolled in the bootcamp.
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} LeaderboardResponse "Leaderboard entries with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid bootcamp ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not enrolled in bootcamp"
// @Failure 404 {object} map[string]any "Not found - bootcamp does not exist"
// @Router /v1/bootcamps/{bootcampId}/leaderboard [get]
func (h *Handler) GetBootcampLeaderboard(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	// Parse pagination parameters
	page := ParsePage((*c).QueryParam("page"))
	limit := ParseLimit((*c).QueryParam("limit"), 20, 100)

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is enrolled in bootcamp
	_, err = h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Get leaderboard entries
	entries, total, err := h.service.GetBootcampLeaderboard(c.Request().Context(), bootcampID, page, limit)
	if err != nil {
		switch err.Error() {
		case "BOOTCAMP_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "BOOTCAMP_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "LEADERBOARD_RETRIEVED", LeaderboardResponse{
		Success: true,
		Data:    entries,
		Meta: &OffsetPagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil)
}

// GetLeaderboardEntry godoc
// @Summary Get leaderboard entry
// @Description Retrieve a specific leaderboard entry by enrollment ID. Mentees can only view their own entry.
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param enrollmentId path string true "Bootcamp Enrollment ID (UUID)"
// @Success 200 {object} LeaderboardEntryResponse "Leaderboard entry details"
// @Failure 400 {object} map[string]any "Bad request - invalid ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - access denied"
// @Failure 404 {object} map[string]any "Not found - entry does not exist"
// @Router /v1/bootcamps/{bootcampId}/leaderboard/{enrollmentId} [get]
func (h *Handler) GetLeaderboardEntry(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	enrollmentID, err := utils.StringToUUID((*c).Param("enrollmentId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ENROLLMENT_ID", nil, nil)
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get member ID for access control
	memberID, err := h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Get leaderboard entry with access control
	entry, err := h.service.GetLeaderboardEntry(c.Request().Context(), bootcampID, enrollmentID, claims.Role, memberID)
	if err != nil {
		switch err.Error() {
		case "ENTRY_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "ENTRY_NOT_FOUND", nil, nil)
		case "ACCESS_DENIED":
			return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ACCESS_DENIED", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "ENTRY_RETRIEVED", entry, nil)
}

// Poll Handlers

// CreatePoll godoc
// @Summary Create a poll
// @Description Create a difficulty poll for a problem in a bootcamp (mentor/admin only). Supports idempotency via Idempotency-Key header.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param body body CreatePollRequest true "Poll details"
// @Param Idempotency-Key header string false "Idempotency key for safe retries"
// @Success 201 {object} PollResponse "Poll created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid problem ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor/admin role required"
// @Failure 404 {object} map[string]any "Not found - problem does not exist"
// @Router /v1/bootcamps/{bootcampId}/polls [post]
func (h *Handler) CreatePoll(c *echo.Context) error {
	var body CreatePollRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}

	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentor/admin
	if claims.Role == "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "MENTOR_ADMIN_ROLE_REQUIRED", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is enrolled in bootcamp
	_, err = h.service.GetMemberIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Create poll
	poll, err := h.service.CreatePoll(c.Request().Context(), bootcampID, body, userID)
	if err != nil {
		switch err.Error() {
		case "INVALID_PROBLEM_ID":
			return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PROBLEM_ID", nil, nil)
		case "PROBLEM_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "PROBLEM_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusCreated, "SUCCESS", "POLL_CREATED", poll, nil)
}

// ListPolls godoc
// @Summary List polls
// @Description List polls for a bootcamp with optional problem filtering. Includes user's vote if they have voted. User must be enrolled in bootcamp.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param problemId query string false "Filter by problem ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} PollListResponse "List of polls with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid bootcamp ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not enrolled in bootcamp"
// @Router /v1/bootcamps/{bootcampId}/polls [get]
func (h *Handler) ListPolls(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	// Parse query parameters
	problemIDStr := (*c).QueryParam("problemId")
	page := ParsePage((*c).QueryParam("page"))
	limit := ParseLimit((*c).QueryParam("limit"), 20, 100)

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is enrolled in bootcamp and get enrollment ID
	enrollmentID, err := h.service.GetEnrollmentIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// List polls
	polls, total, err := h.service.ListPolls(c.Request().Context(), bootcampID, problemIDStr, enrollmentID, page, limit)
	if err != nil {
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "POLLS_RETRIEVED", PollListResponse{
		Success: true,
		Data:    polls,
		Meta: &OffsetPagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil)
}

// GetPoll godoc
// @Summary Get poll details
// @Description Retrieve full details of a specific poll including user's vote state. User must be enrolled in bootcamp.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param pollId path string true "Poll ID (UUID)"
// @Success 200 {object} PollResponse "Poll details"
// @Failure 400 {object} map[string]any "Bad request - invalid poll ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - not enrolled in bootcamp"
// @Failure 404 {object} map[string]any "Not found - poll does not exist"
// @Router /v1/bootcamps/{bootcampId}/polls/{pollId} [get]
func (h *Handler) GetPoll(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	pollID, err := utils.StringToUUID((*c).Param("pollId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_POLL_ID", nil, nil)
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Verify user is enrolled in bootcamp and get enrollment ID
	enrollmentID, err := h.service.GetEnrollmentIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Get poll
	poll, err := h.service.GetPoll(c.Request().Context(), bootcampID, pollID, enrollmentID)
	if err != nil {
		switch err.Error() {
		case "POLL_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "POLL_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "POLL_RETRIEVED", poll, nil)
}

// VotePoll godoc
// @Summary Vote on a poll
// @Description Cast or update a vote on a poll (mentee only). Uses PUT method for idempotent vote creation/update. Returns 201 for first vote, 200 for updates.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param pollId path string true "Poll ID (UUID)"
// @Param body body VotePollRequest true "Vote details"
// @Success 200 {object} VoteResponse "Vote updated successfully"
// @Success 201 {object} VoteResponse "Vote created successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid poll ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - only mentees can vote"
// @Failure 404 {object} map[string]any "Not found - poll does not exist"
// @Router /v1/bootcamps/{bootcampId}/polls/{pollId}/vote [put]
func (h *Handler) VotePoll(c *echo.Context) error {
	var body VotePollRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}

	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentee
	if claims.Role != "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ONLY_MENTEES_CAN_VOTE", nil, nil)
	}

	bootcampID, err := utils.StringToUUID((*c).Param("bootcampId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BOOTCAMP_ID", nil, nil)
	}

	pollID, err := utils.StringToUUID((*c).Param("pollId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_POLL_ID", nil, nil)
	}

	// Get user ID
	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	// Get voter's enrollment ID
	voterEnrollmentID, err := h.service.GetEnrollmentIDByUserAndBootcamp(c.Request().Context(), userID, bootcampID)
	if err != nil {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "NOT_ENROLLED_IN_BOOTCAMP", nil, nil)
	}

	// Validate poll belongs to bootcamp
	poll, err := h.service.GetPoll(c.Request().Context(), bootcampID, pollID, voterEnrollmentID)
	if err != nil {
		switch err.Error() {
		case "POLL_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "POLL_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	// Cast vote
	vote, isNew, err := h.service.VotePoll(c.Request().Context(), poll.Data.ID, voterEnrollmentID, body.Vote)
	if err != nil {
		switch err.Error() {
		case "POLL_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "POLL_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	statusCode := http.StatusOK
	message := "VOTE_UPDATED"
	if isNew {
		statusCode = http.StatusCreated
		message = "VOTE_CREATED"
	}

	return response.NewResponse(c, statusCode, "SUCCESS", message, vote, nil)
}

// GetPollResults godoc
// @Summary Get poll results
// @Description Retrieve aggregated poll results with vote counts and percentages (mentor/admin/super_admin only). Mentees cannot access results.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param pollId path string true "Poll ID (UUID)"
// @Success 200 {object} PollResultsResponse "Aggregated poll results"
// @Failure 400 {object} map[string]any "Bad request - invalid poll ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor/admin/super_admin role required"
// @Failure 404 {object} map[string]any "Not found - poll does not exist"
// @Router /v1/bootcamps/{bootcampId}/polls/{pollId}/results [get]
func (h *Handler) GetPollResults(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentor/admin/super_admin
	if claims.Role == "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "MENTEES_CANNOT_ACCESS_RESULTS", nil, nil)
	}

	pollID, err := utils.StringToUUID((*c).Param("pollId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_POLL_ID", nil, nil)
	}

	// Get poll results
	results, err := h.service.GetPollResults(c.Request().Context(), pollID)
	if err != nil {
		switch err.Error() {
		case "POLL_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "POLL_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "RESULTS_RETRIEVED", results, nil)
}

// GetPollVotes godoc
// @Summary Get individual poll votes
// @Description Retrieve individual vote records with optional filtering by vote value (mentor/admin/super_admin only). Includes voter enrollment ID but not internal user identifiers.
// @Tags Polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bootcampId path string true "Bootcamp ID (UUID)"
// @Param pollId path string true "Poll ID (UUID)"
// @Param vote query string false "Filter by vote value (easy, medium, hard)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} PollVotesResponse "List of individual votes with pagination"
// @Failure 400 {object} map[string]any "Bad request - invalid poll ID"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]any "Forbidden - mentor/admin/super_admin role required"
// @Failure 404 {object} map[string]any "Not found - poll does not exist"
// @Router /v1/bootcamps/{bootcampId}/polls/{pollId}/votes [get]
func (h *Handler) GetPollVotes(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	// Validate user is mentor/admin/super_admin
	if claims.Role == "mentee" {
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "MENTEES_CANNOT_ACCESS_VOTES", nil, nil)
	}

	pollID, err := utils.StringToUUID((*c).Param("pollId"))
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_POLL_ID", nil, nil)
	}

	// Parse query parameters
	voteFilter := (*c).QueryParam("vote")
	page := ParsePage((*c).QueryParam("page"))
	limit := ParseLimit((*c).QueryParam("limit"), 20, 100)

	// Get poll votes
	votes, total, err := h.service.GetPollVotes(c.Request().Context(), pollID, voteFilter, page, limit)
	if err != nil {
		switch err.Error() {
		case "POLL_NOT_FOUND":
			return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "POLL_NOT_FOUND", nil, nil)
		default:
			return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, nil)
		}
	}

	return response.NewResponse(c, http.StatusOK, "SUCCESS", "VOTES_RETRIEVED", PollVotesResponse{
		Success: true,
		Data:    votes,
		Meta: &OffsetPagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil)
}
