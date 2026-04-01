-- name: CreateAssignmentGroup :one
INSERT INTO assignment_groups (
    bootcamp_id, created_by, title, description, deadline_days
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetAssignmentGroup :one
SELECT * FROM assignment_groups
WHERE id = $1 LIMIT 1;

-- name: UpdateAssignmentGroup :one
UPDATE assignment_groups
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    deadline_days = COALESCE(sqlc.narg('deadline_days'), deadline_days),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListAssignmentGroupsByBootcamp :many
SELECT * FROM assignment_groups
WHERE bootcamp_id = $1
  AND (sqlc.narg('created_by')::uuid IS NULL OR created_by = sqlc.narg('created_by')::uuid)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountAssignmentGroupsByBootcamp :one
SELECT COUNT(*) FROM assignment_groups
WHERE bootcamp_id = $1
  AND (sqlc.narg('created_by')::uuid IS NULL OR created_by = sqlc.narg('created_by')::uuid);

-- name: AddProblemToAssignmentGroup :exec
INSERT INTO assignment_group_problems (
    assignment_group_id, problem_id, position
) VALUES (
    $1, $2, $3
)
ON CONFLICT (assignment_group_id, problem_id) DO UPDATE SET position = EXCLUDED.position;

-- name: RemoveProblemFromAssignmentGroup :exec
DELETE FROM assignment_group_problems
WHERE assignment_group_id = $1 AND problem_id = $2;

-- name: ClearAssignmentGroupProblems :exec
DELETE FROM assignment_group_problems
WHERE assignment_group_id = $1;

-- name: ListAssignmentGroupProblems :many
SELECT p.*, agp.position 
FROM problems p
JOIN assignment_group_problems agp ON p.id = agp.problem_id
WHERE agp.assignment_group_id = $1
ORDER BY agp.position ASC;

-- Assignment Instances

-- name: AssignGroupToMentee :one
INSERT INTO assignments (
    assignment_group_id, bootcamp_enrollment_id, assigned_by, deadline_at, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetAssignment :one
SELECT * FROM assignments
WHERE id = $1 AND archived_at IS NULL LIMIT 1;

-- name: GetAssignmentWithGroup :one
SELECT a.*, ag.title as group_title, ag.description as group_description
FROM assignments a
JOIN assignment_groups ag ON a.assignment_group_id = ag.id
WHERE a.id = $1 AND a.archived_at IS NULL LIMIT 1;

-- name: ListAssignmentsByMentee :many
SELECT a.*, ag.title as group_title 
FROM assignments a
JOIN assignment_groups ag ON a.assignment_group_id = ag.id
WHERE a.bootcamp_enrollment_id = $1 AND a.archived_at IS NULL
ORDER BY a.deadline_at ASC;

-- name: ListAssignments :many
SELECT a.*, ag.title as group_title
FROM assignments a
JOIN assignment_groups ag ON a.assignment_group_id = ag.id
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE be.bootcamp_id = $1
  AND (sqlc.narg('assignment_group_id')::uuid IS NULL OR a.assignment_group_id = sqlc.narg('assignment_group_id')::uuid)
  AND (sqlc.narg('status')::assignment_status IS NULL OR a.status = sqlc.narg('status')::assignment_status)
  AND a.archived_at IS NULL
ORDER BY a.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountAssignments :one
SELECT COUNT(*)
FROM assignments a
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE be.bootcamp_id = $1
  AND (sqlc.narg('assignment_group_id')::uuid IS NULL OR a.assignment_group_id = sqlc.narg('assignment_group_id')::uuid)
  AND (sqlc.narg('status')::assignment_status IS NULL OR a.status = sqlc.narg('status')::assignment_status)
  AND a.archived_at IS NULL;

-- name: CheckDuplicateActiveAssignment :one
SELECT COUNT(*) FROM assignments
WHERE assignment_group_id = $1 
  AND bootcamp_enrollment_id = $2 
  AND status = 'active'
  AND archived_at IS NULL;

-- name: GetEnrollmentBootcamp :one
SELECT be.bootcamp_id, b.is_active
FROM bootcamp_enrollments be
JOIN bootcamps b ON be.bootcamp_id = b.id
WHERE be.id = $1;

-- name: UpdateAssignmentStatus :one
UPDATE assignments
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateAssignmentDeadline :one
UPDATE assignments
SET deadline_at = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ArchiveAssignment :exec
UPDATE assignments
SET archived_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Assignment Problems Progress

-- name: InitializeAssignmentProblem :one
INSERT INTO assignment_problems (
    assignment_id, problem_id, status
) VALUES (
    $1, $2, 'pending'
)
RETURNING *;

-- name: InitializeAssignmentProblems :exec
INSERT INTO assignment_problems (
    assignment_id, problem_id, status
)
SELECT $1, unnest(sqlc.arg('problem_ids')::uuid[]), 'pending';

-- name: UpdateAssignmentProblemProgress :one
UPDATE assignment_problems
SET 
    status = COALESCE(sqlc.narg('status'), status),
    solution_link = COALESCE(sqlc.narg('solution_link'), solution_link),
    notes = COALESCE(sqlc.narg('notes'), notes),
    completed_at = COALESCE(sqlc.narg('completed_at'), completed_at),
    updated_at = CURRENT_TIMESTAMP
WHERE assignment_id = $1 AND problem_id = $2
RETURNING *;

-- name: ListAssignmentProblemsStatus :many
SELECT ap.*, p.title, p.difficulty 
FROM assignment_problems ap
JOIN problems p ON ap.problem_id = p.id
WHERE ap.assignment_id = $1
ORDER BY ap.created_at ASC;

-- name: GetAssignmentProblem :one
SELECT ap.*, p.title, p.difficulty 
FROM assignment_problems ap
JOIN problems p ON ap.problem_id = p.id
WHERE ap.assignment_id = $1 AND ap.problem_id = $2
LIMIT 1;

-- name: GetAssignmentWithEnrollment :one
SELECT a.*, be.organization_member_id
FROM assignments a
JOIN bootcamp_enrollments be ON a.bootcamp_enrollment_id = be.id
WHERE a.id = $1 AND a.archived_at IS NULL
LIMIT 1;

-- name: CountAssignmentsByGroup :one
SELECT COUNT(*) FROM assignments
WHERE assignment_group_id = $1 AND archived_at IS NULL;

-- name: DeleteAssignmentGroup :exec
DELETE FROM assignment_groups
WHERE id = $1;
