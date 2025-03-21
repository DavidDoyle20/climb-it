// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_tokens.sql

package database

import (
	"context"
)

const assignRefreshTokenToUser = `-- name: AssignRefreshTokenToUser :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    datetime('now', '+60 days'),
    NULL
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type AssignRefreshTokenToUserParams struct {
	Token  string
	UserID string
}

func (q *Queries) AssignRefreshTokenToUser(ctx context.Context, arg AssignRefreshTokenToUserParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, assignRefreshTokenToUser, arg.Token, arg.UserID)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const checkAndFetchRefreshToken = `-- name: CheckAndFetchRefreshToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = ?
    AND expires_at > datetime('now')
    AND revoked_at IS NULL
`

func (q *Queries) CheckAndFetchRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, checkAndFetchRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getRefreshTokenFromUser = `-- name: GetRefreshTokenFromUser :one
SELECT refresh_tokens.token, refresh_tokens.created_at, refresh_tokens.updated_at, refresh_tokens.user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at
FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE users.id = ? 
    AND expires_at > datetime('now')
    AND revoked_at IS NULL
`

func (q *Queries) GetRefreshTokenFromUser(ctx context.Context, id string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshTokenFromUser, id)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT users.id, users.created_at, users.updated_at, users.name, users.email, users.hashed_password
FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = ?
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const revokeRefreshTokenFromUser = `-- name: RevokeRefreshTokenFromUser :exec
DELETE FROM refresh_tokens
WHERE user_id = ?
`

func (q *Queries) RevokeRefreshTokenFromUser(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshTokenFromUser, userID)
	return err
}
