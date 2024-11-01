-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, email, password)
VALUES (?, ?, ?, ?);

-- name: UpdateUser :exec
UPDATE users
SET first_name = ?, last_name = ?, email = ?, password = ?
WHERE user_id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = ?;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = ?;

-- name: GetUserByEmail :one 
SELECT * FROM users WHERE email = ?;

