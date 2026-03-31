-- name: CreateOrganization :one
INSERT INTO organizations (
    name, slug, description, status
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetOrganizationById :one
SELECT * FROM organizations
WHERE id = $1 LIMIT 1;

-- name: GetOrganizationBySlug :one
SELECT * FROM organizations
WHERE slug = $1 LIMIT 1;

-- name: ListOrganizations :many
SELECT o.* FROM organizations o
JOIN organization_members om ON o.id = om.organization_id
WHERE om.user_id = $1
ORDER BY o.created_at DESC;

-- name: UpdateOrganization :one
UPDATE organizations
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetPendingOrganizations :many
SELECT * FROM organizations
WHERE status = 'pending_approval'
ORDER BY created_at ASC;

-- Member management

-- name: AddOrganizationMember :one
INSERT INTO organization_members (
    organization_id, user_id, role
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetOrganizationMember :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND user_id = $2 LIMIT 1;

-- name: ListOrganizationMembers :many
SELECT om.*, u.name, u.email, u.avatar_url 
FROM organization_members om
JOIN users u ON om.user_id = u.id
WHERE om.organization_id = $1
ORDER BY om.joined_at ASC;

-- name: UpdateOrganizationMemberRole :one
UPDATE organization_members
SET role = $3
WHERE organization_id = $1 AND user_id = $2
RETURNING *;

-- name: RemoveOrganizationMember :exec
DELETE FROM organization_members
WHERE organization_id = $1 AND user_id = $2;
