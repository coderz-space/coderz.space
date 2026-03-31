-- name: UpsertLeaderboardEntry :one
INSERT INTO leaderboard_entries (
    bootcamp_id, bootcamp_enrollment_id, problems_completed, problems_attempted, 
    completion_rate, streak_days, score, rank, calculated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP
)
ON CONFLICT (bootcamp_id, bootcamp_enrollment_id) DO UPDATE SET
    problems_completed = EXCLUDED.problems_completed,
    problems_attempted = EXCLUDED.problems_attempted,
    completion_rate = EXCLUDED.completion_rate,
    streak_days = EXCLUDED.streak_days,
    score = EXCLUDED.score,
    rank = EXCLUDED.rank,
    calculated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetLeaderboardByBootcamp :many
SELECT le.*, u.name, u.avatar_url
FROM leaderboard_entries le
JOIN bootcamp_enrollments be ON le.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
JOIN users u ON om.user_id = u.id
WHERE le.bootcamp_id = $1
ORDER BY le.rank ASC;

-- Polls

-- name: CreatePoll :one
INSERT INTO polls (
    bootcamp_id, problem_id, question, created_by
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetPoll :one
SELECT * FROM polls
WHERE id = $1 LIMIT 1;

-- name: ListPollsByBootcamp :many
SELECT p.*, prob.title as problem_title
FROM polls p
JOIN problems prob ON p.problem_id = prob.id
WHERE p.bootcamp_id = $1
ORDER BY p.created_at DESC;

-- name: CastPollVote :one
INSERT INTO poll_votes (
    poll_id, voter_id, vote
) VALUES (
    $1, $2, $3
)
ON CONFLICT (poll_id, voter_id) DO UPDATE SET vote = EXCLUDED.vote
RETURNING *;

-- name: GetPollResults :many
SELECT vote, COUNT(*) as vote_count
FROM poll_votes
WHERE poll_id = $1
GROUP BY vote;

-- name: GetUserVoteForPoll :one
SELECT * FROM poll_votes
WHERE poll_id = $1 AND voter_id = $2
LIMIT 1;

-- name: ListPollVotesByPoll :many
SELECT pv.*, u.name as voter_name
FROM poll_votes pv
JOIN bootcamp_enrollments be ON pv.voter_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
JOIN users u ON om.user_id = u.id
WHERE pv.poll_id = $1
  AND ($2::text IS NULL OR pv.vote = $2)
ORDER BY pv.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountPollVotesByPoll :one
SELECT COUNT(*) FROM poll_votes
WHERE poll_id = $1
  AND ($2::text IS NULL OR vote = $2);

-- name: CheckVoteExists :one
SELECT EXISTS(
    SELECT 1 FROM poll_votes
    WHERE poll_id = $1 AND voter_id = $2
) as vote_exists;
