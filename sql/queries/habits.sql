-- name: CreateHabitForUser :one
INSERT INTO habits (id, created_at, updated_at, name, user_id, start_date, end_date)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: RemoveHabit :exec
DELETE FROM habits
WHERE id = ?;

-- name: GetHabit :one
SELECT * FROM habits
WHERE id = ?;

-- name: GetUserHabits :many
SELECT * FROM habits
WHERE user_id = ?;