-- name: CreateBootcamp :one
INSERT INTO bootcamps (
    organization_id, created_by, name, description, start_date, end_date, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetBootcamp :one
SELECT * FROM bootcamps
WHERE id = $1 AND archived_at IS NULL LIMIT 1;

-- name: ListBootcampsByOrg :many
SELECT * FROM bootcamps
WHERE organization_id = $1 AND archived_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateBootcamp :one
UPDATE bootcamps
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ArchiveBootcamp :exec
UPDATE bootcamps
SET archived_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Enrollment

-- name: EnrollInBootcamp :one
INSERT INTO bootcamp_enrollments (
    bootcamp_id, organization_member_id, role, status
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetEnrollment :one
SELECT * FROM bootcamp_enrollments
WHERE id = $1 LIMIT 1;

-- name: GetEnrollmentByMember :one
SELECT * FROM bootcamp_enrollments
WHERE bootcamp_id = $1 AND organization_member_id = $2 LIMIT 1;

-- name: ListBootcampEnrollments :many
SELECT be.*, u.name, u.email, u.avatar_url, om.role as org_role
FROM bootcamp_enrollments be
JOIN organization_members om ON be.organization_member_id = om.id
JOIN users u ON om.user_id = u.id
WHERE be.bootcamp_id = $1
ORDER BY be.enrolled_at ASC;

-- name: UpdateEnrollmentRole :one
UPDATE bootcamp_enrollments
SET role = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEnrollmentStatus :one
UPDATE bootcamp_enrollments
SET status = $2
WHERE id = $1
RETURNING *;

-- name: RemoveEnrollment :exec
DELETE FROM bootcamp_enrollments
WHERE id = $1;
