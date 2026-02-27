-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password, active, email_verification_hash, token_expires_at
) VALUES (
    $1, $2, $3, false, $4, $5
)
RETURNING *;

-- name: UpdateUserProfile :execrows
UPDATE users
SET
    username = $2,
    email = $3
WHERE id = $1;

-- name: UpdateUserPassword :execrows
UPDATE users
SET
    password = $2
WHERE id = $1;

-- name: UpdateUserActive :execrows
UPDATE users
SET
    active = $2
WHERE id = $1;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;
