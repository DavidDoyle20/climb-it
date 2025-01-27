-- name: AssignRefreshTokenToUser :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    ?,
    time('now'),
    time('now'),
    ?,
    time('now', '+60 days'),
    NULL
)
RETURNING *;

-- name: RevokeRefreshTokenFromUser :exec
UPDATE refresh_tokens
SET revoked_at = time('now'), updated_at = time('now')
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
    AND expires_at > time('now')
    AND revoked_at IS NULL;

