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
ON CONFLICT (organization_id, name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: ListTagsByOrg :many
SELECT * FROM tags
WHERE organization_id = $1
ORDER BY name ASC;

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

-- Resources

-- name: AddProblemResource :one
INSERT INTO problem_resources (
    problem_id, title, url
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ListProblemResources :many
SELECT * FROM problem_resources
WHERE problem_id = $1
ORDER BY created_at ASC;

-- name: DeleteProblemResource :exec
DELETE FROM problem_resources
WHERE id = $1;
