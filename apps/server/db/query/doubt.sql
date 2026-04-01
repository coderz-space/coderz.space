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

-- name: GetDoubtWithDetails :one
SELECT 
    d.*,
    u_raised.name as raised_by_name,
    u_raised.email as raised_by_email,
    u_resolved.name as resolved_by_name
FROM doubts d
JOIN organization_members om_raised ON d.raised_by = om_raised.id
JOIN users u_raised ON om_raised.user_id = u_raised.id
LEFT JOIN organization_members om_resolved ON d.resolved_by = om_resolved.id
LEFT JOIN users u_resolved ON om_resolved.user_id = u_resolved.id
WHERE d.id = $1 LIMIT 1;

-- name: ListDoubtsByAssignmentProblem :many
SELECT d.*, u.name as raised_by_name 
FROM doubts d
JOIN organization_members om ON d.raised_by = om.id
JOIN users u ON om.user_id = u.id
WHERE d.assignment_problem_id = $1
ORDER BY d.created_at DESC;

-- name: ListDoubtsByMentee :many
SELECT 
    d.*,
    u_raised.name as raised_by_name,
    u_raised.email as raised_by_email,
    u_resolved.name as resolved_by_name
FROM doubts d
JOIN organization_members om_raised ON d.raised_by = om_raised.id
JOIN users u_raised ON om_raised.user_id = u_raised.id
LEFT JOIN organization_members om_resolved ON d.resolved_by = om_resolved.id
LEFT JOIN users u_resolved ON om_resolved.user_id = u_resolved.id
WHERE d.raised_by = $1
    AND ($2::boolean IS NULL OR d.resolved = $2)
ORDER BY d.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountDoubtsByMentee :one
SELECT COUNT(*) FROM doubts
WHERE raised_by = $1
    AND ($2::boolean IS NULL OR resolved = $2);

-- name: ListDoubtsByBootcamp :many
SELECT 
    d.*,
    u_raised.name as raised_by_name,
    u_raised.email as raised_by_email,
    u_resolved.name as resolved_by_name
FROM doubts d
JOIN assignment_problems ap ON d.assignment_problem_id = ap.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om_raised ON d.raised_by = om_raised.id
JOIN users u_raised ON om_raised.user_id = u_raised.id
LEFT JOIN organization_members om_resolved ON d.resolved_by = om_resolved.id
LEFT JOIN users u_resolved ON om_resolved.user_id = u_resolved.id
WHERE be.bootcamp_id = $1
    AND ($2::uuid IS NULL OR d.assignment_problem_id = $2)
    AND ($3::boolean IS NULL OR d.resolved = $3)
ORDER BY d.created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountDoubtsByBootcamp :one
SELECT COUNT(*) FROM doubts d
JOIN assignment_problems ap ON d.assignment_problem_id = ap.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE be.bootcamp_id = $1
    AND ($2::uuid IS NULL OR d.assignment_problem_id = $2)
    AND ($3::boolean IS NULL OR d.resolved = $3);

-- name: ListDoubtsCursor :many
SELECT 
    d.*,
    u_raised.name as raised_by_name,
    u_raised.email as raised_by_email,
    u_resolved.name as resolved_by_name
FROM doubts d
JOIN assignment_problems ap ON d.assignment_problem_id = ap.id
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
JOIN organization_members om_raised ON d.raised_by = om_raised.id
JOIN users u_raised ON om_raised.user_id = u_raised.id
LEFT JOIN organization_members om_resolved ON d.resolved_by = om_resolved.id
LEFT JOIN users u_resolved ON om_resolved.user_id = u_resolved.id
WHERE be.bootcamp_id = $1
    AND ($2::uuid IS NULL OR d.assignment_problem_id = $2)
    AND ($3::boolean IS NULL OR d.resolved = $3)
    AND ($4::uuid IS NULL OR d.id < $4)
ORDER BY d.created_at DESC, d.id DESC
LIMIT $5;

-- name: ListDoubtsByMenteeCursor :many
SELECT 
    d.*,
    u_raised.name as raised_by_name,
    u_raised.email as raised_by_email,
    u_resolved.name as resolved_by_name
FROM doubts d
JOIN organization_members om_raised ON d.raised_by = om_raised.id
JOIN users u_raised ON om_raised.user_id = u_raised.id
LEFT JOIN organization_members om_resolved ON d.resolved_by = om_resolved.id
LEFT JOIN users u_resolved ON om_resolved.user_id = u_resolved.id
WHERE d.raised_by = $1
    AND ($2::boolean IS NULL OR d.resolved = $2)
    AND ($3::uuid IS NULL OR d.id < $3)
ORDER BY d.created_at DESC, d.id DESC
LIMIT $4;

-- name: ResolveDoubt :one
UPDATE doubts
SET 
    resolved = TRUE,
    resolved_by = $2,
    resolved_at = CURRENT_TIMESTAMP,
    resolution_note = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteDoubt :exec
DELETE FROM doubts
WHERE id = $1;

-- name: GetAssignmentProblemDetails :one
SELECT 
    ap.id,
    ap.assignment_id,
    a.bootcamp_enrollment_id,
    be.bootcamp_id,
    be.organization_member_id
FROM assignment_problems ap
JOIN assignments a ON ap.assignment_id = a.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE ap.id = $1;

-- name: ValidateAssignmentProblemOwnership :one
SELECT EXISTS(
    SELECT 1 FROM assignment_problems ap
    JOIN assignments a ON ap.assignment_id = a.id
    JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
    WHERE ap.id = $1 AND be.organization_member_id = $2
) as is_owner;

-- name: GetEnrollmentByMemberID :one
SELECT be.* FROM bootcamp_enrollments be
WHERE be.organization_member_id = $1 AND be.bootcamp_id = $2
LIMIT 1;

-- name: GetMemberIDByUserID :one
SELECT om.id FROM organization_members om
JOIN bootcamp_enrollments be ON om.id = be.organization_member_id
WHERE om.user_id = $1 AND be.bootcamp_id = $2
LIMIT 1;

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
-- name: ValidateDoubtResolverOrg :one
SELECT EXISTS(
    SELECT 1 FROM doubts d
    JOIN assignment_problems ap ON d.assignment_problem_id = ap.id
    JOIN assignments a ON ap.assignment_id = a.id
    JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
    JOIN organization_members mentee_om ON be.organization_member_id = mentee_om.id
    JOIN organization_members resolver_om ON resolver_om.id = sqlc.arg('resolver_member_id')
    WHERE d.id = sqlc.arg('doubt_id') AND mentee_om.organization_id = resolver_om.organization_id
) as is_same_org;
