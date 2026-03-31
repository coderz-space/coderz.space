-- name: CreateUser :one
INSERT INTO users (
    name, email, password_hash, google_id, avatar_url, role
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByGoogleId :one
SELECT * FROM users
WHERE google_id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    email_verified = COALESCE(sqlc.narg('email_verified'), email_verified),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, token_hash, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1 LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: ClearExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (
    user_id, token_hash, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens
WHERE token_hash = $1;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: DeleteUserPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE user_id = $1;
