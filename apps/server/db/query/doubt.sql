-- name: CreateDoubt :one
INSERT INTO doubts (
    assignment_problem_id, raised_by, message
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetDoubt :one
SELECT * FROM doubts
WHERE id = $1 LIMIT 1;

-- name: ListDoubtsByAssignmentProblem :many
SELECT d.*, u.name as raised_by_name 
FROM doubts d
JOIN organization_members om ON d.raised_by = om.id
JOIN users u ON om.user_id = u.id
WHERE d.assignment_problem_id = $1
ORDER BY d.created_at DESC;

-- name: ResolveDoubt :one
UPDATE doubts
SET 
    resolved = TRUE,
    resolved_by = $2,
    resolved_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListPendingDoubtsByBootcamp :many
SELECT d.*, p.title as problem_title, u.name as mentee_name
FROM doubts d
JOIN assignment_problems ap ON d.assignment_problem_id = ap.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN problems p ON ap.problem_id = p.id
JOIN organization_members om ON d.raised_by = om.id
JOIN users u ON om.user_id = u.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE be.bootcamp_id = $1 AND d.resolved = FALSE
ORDER BY d.created_at ASC;
