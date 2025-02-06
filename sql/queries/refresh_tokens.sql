-- name: AssignRefreshTokenToUser :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    datetime('now', '+60 days'),
    NULL
)
RETURNING *;

-- name: RevokeRefreshTokenFromUser :exec
DELETE FROM refresh_tokens
WHERE user_id = ?;

-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = ?;

-- name: CheckAndFetchRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = ?
    AND expires_at > datetime('now')
    AND revoked_at IS NULL;

-- name: GetRefreshTokenFromUser :one
SELECT refresh_tokens.*
FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE users.id = ? 
    AND expires_at > datetime('now')
    AND revoked_at IS NULL;

