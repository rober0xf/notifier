-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT
    1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT
    1;

-- name: GetUserByGoogleID :one
SELECT
    *
FROM
    users
WHERE
    google_id = $1
LIMIT
    1;

-- name: GetAllUsers :many
SELECT
    *
FROM
    users;

-- name: CreateUser :one
INSERT INTO
    users (username, email, PASSWORD, active)
VALUES
    ($1, $2, $3, FALSE)
RETURNING
    *;

-- name: CreateOAuthUser :one
INSERT INTO
    users (username, email, name, google_id, active)
VALUES
    ($1, $2, $3, $4, TRUE)
RETURNING
    *;

-- name: UpdateUserProfile :execrows
UPDATE
    users
SET
    username = $2,
    email = $3
WHERE
    id = $1;

-- name: UpdateUserPassword :execrows
UPDATE
    users
SET
    PASSWORD = $2
WHERE
    id = $1;

-- name: UpdateUserActive :execrows
UPDATE
    users
SET
    active = $2
WHERE
    id = $1;

-- name: UpdateUserGoogleID :execrows
UPDATE
    users
SET
    google_id = $2
WHERE
    id = $1
    AND google_id IS NULL;

-- name: DeleteUser :execrows
DELETE FROM
    users
WHERE
    id = $1;

-- user tokens queries
-- name: CreateUserToken :one
INSERT INTO
    user_tokens (user_id, token_hash, purpose, expires_at)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: VerifyToken :one
UPDATE
    user_tokens
SET
    used = TRUE
WHERE
    token_hash = $1
    AND purpose = $2
    AND used = false
    AND expires_at > NOW()
RETURNING
    *;

-- name: GetValidTokenByHash :one
SELECT
    *
FROM
    user_tokens
WHERE
    token_hash = $1
    and purpose = $2
    AND used = false
    AND expires_at > NOW();

-- name: DeleteOldTokens :execrows
DELETE FROM
    user_tokens
WHERE
    used = TRUE
    OR expires_at < NOW();
