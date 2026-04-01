-- name: CreateProblem :one
INSERT INTO problems (
    organization_id, created_by, title, description, difficulty, external_link
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetProblem :one
SELECT * FROM problems
WHERE id = $1 AND archived_at IS NULL LIMIT 1;

-- name: ListProblemsByOrg :many
SELECT * FROM problems
WHERE organization_id = $1 AND archived_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateProblem :one
UPDATE problems
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    difficulty = COALESCE(sqlc.narg('difficulty'), difficulty),
    external_link = COALESCE(sqlc.narg('external_link'), external_link),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ArchiveProblem :exec
UPDATE problems
SET archived_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Tags

-- name: CreateTag :one
INSERT INTO tags (
    organization_id, created_by, name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTag :one
SELECT * FROM tags
WHERE id = $1 LIMIT 1;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE organization_id = $1 AND name = $2 LIMIT 1;

-- name: ListTagsByOrg :many
SELECT * FROM tags
WHERE organization_id = $1
ORDER BY name ASC;

-- name: SearchTagsByName :many
SELECT * FROM tags
WHERE organization_id = $1 AND name ILIKE '%' || sqlc.arg('name')::text || '%'
ORDER BY name ASC;

-- name: UpdateTag :one
UPDATE tags
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1;

-- name: CountTagUsage :one
SELECT COUNT(*) FROM problem_tags
WHERE tag_id = $1;

-- name: AddTagToProblem :exec
INSERT INTO problem_tags (problem_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromProblem :exec
DELETE FROM problem_tags
WHERE problem_id = $1 AND tag_id = $2;

-- name: ListProblemTags :many
SELECT t.* FROM tags t
JOIN problem_tags pt ON t.id = pt.tag_id
WHERE pt.problem_id = $1;

-- name: GetTagsByIDs :many
SELECT * FROM tags
WHERE id = ANY($1::uuid[]);

-- Resources

-- name: AddProblemResource :one
INSERT INTO problem_resources (
    problem_id, title, url
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetProblemResource :one
SELECT * FROM problem_resources
WHERE id = $1 LIMIT 1;

-- name: ListProblemResources :many
SELECT * FROM problem_resources
WHERE problem_id = $1
ORDER BY created_at ASC;

-- name: UpdateProblemResource :one
UPDATE problem_resources
SET 
    title = COALESCE(sqlc.narg('title'), title),
    url = COALESCE(sqlc.narg('url'), url)
WHERE id = $1
RETURNING *;

-- name: DeleteProblemResource :exec
DELETE FROM problem_resources
WHERE id = $1;

-- Super Admin Queries

-- name: ListAllProblems :many
SELECT p.*, o.name as organization_name, o.slug as organization_slug
FROM problems p
JOIN organizations o ON p.organization_id = o.id
WHERE p.archived_at IS NULL
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllProblems :one
SELECT COUNT(*) FROM problems
WHERE archived_at IS NULL;

-- name: CountProblemAssignments :one
SELECT COUNT(*) FROM assignment_group_problems
WHERE problem_id = $1;
