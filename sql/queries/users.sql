-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, hashed_password)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
--


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;
--