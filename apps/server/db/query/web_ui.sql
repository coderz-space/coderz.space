-- Web UI Specific Queries for Coderz Dashboard

-- name: WebListPendingRequests :many
SELECT be.id, u.id as user_id, u.name as first_name, u.email, be.enrolled_at as signed_up_at, be.status
FROM bootcamp_enrollments be
JOIN organization_members om ON be.organization_member_id = om.id
JOIN users u ON om.user_id = u.id
WHERE be.status = 'pending'
ORDER BY be.enrolled_at DESC;

-- name: WebListMenteeQuestions :many
SELECT 
    ap.problem_id as id,
    p.title,
    p.description,
    p.difficulty,
    a.status as assignment_status,
    ap.status as progress_status,
    a.assigned_at,
    ap.completed_at,
    ap.solution_link as solution_url,
    ap.notes as solution,
    '' as resources
FROM assignment_problems ap
JOIN problems p ON ap.problem_id = p.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
WHERE om.user_id = $1
ORDER BY ap.created_at DESC;

-- name: WebGetMenteeQuestion :one
SELECT 
    ap.problem_id as id,
    p.title,
    p.description,
    p.difficulty,
    a.status as assignment_status,
    ap.status as progress_status,
    a.assigned_at,
    ap.completed_at,
    ap.solution_link as solution_url,
    ap.notes as solution
FROM assignment_problems ap
JOIN problems p ON ap.problem_id = p.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
WHERE om.user_id = $1 AND p.id = $2
LIMIT 1;

-- name: WebUpdateQuestionProgress :exec
UPDATE assignment_problems
SET status = $3, updated_at = CURRENT_TIMESTAMP,
    completed_at = CASE WHEN $3 = 'completed' THEN CURRENT_TIMESTAMP ELSE completed_at END
FROM assignments a
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
WHERE assignment_problems.assignment_id = a.id 
  AND om.user_id = $1 
  AND assignment_problems.problem_id = $2;

-- name: WebUpdateQuestionDetails :exec
UPDATE assignment_problems
SET notes = COALESCE(sqlc.narg('notes'), notes),
    solution_link = COALESCE(sqlc.narg('solution_link'), solution_link),
    updated_at = CURRENT_TIMESTAMP
FROM assignments a
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
WHERE assignment_problems.assignment_id = a.id 
  AND om.user_id = $1 
  AND assignment_problems.problem_id = $2;

-- name: WebGetProfileStats :one
SELECT 
    u.name as first_name,
    u.email,
    u.created_at as joined_at,
    COALESCE(le.problems_completed, 0) as solved
FROM users u
LEFT JOIN organization_members om ON om.user_id = u.id
LEFT JOIN bootcamp_enrollments be ON be.organization_member_id = om.id AND be.status = 'active'
LEFT JOIN leaderboard_entries le ON le.bootcamp_enrollment_id = be.id
WHERE u.id = $1
LIMIT 1;

-- name: WebGetLeaderboard :many
SELECT 
    u.id,
    u.name as first_name,
    '' as last_name,
    le.problems_completed as solved
FROM leaderboard_entries le
JOIN bootcamp_enrollments be ON le.bootcamp_enrollment_id = be.id
JOIN organization_members om ON be.organization_member_id = om.id
JOIN users u ON om.user_id = u.id
ORDER BY le.problems_completed DESC, le.calculated_at DESC
LIMIT 100;
