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
  users.*
FROM
  users
ORDER BY
  id
LIMIT
  $1
OFFSET
  $2;

-- name: CreateUser :one
INSERT INTO
  users (username, email, password_hash, is_active)
VALUES
  ($1, $2, $3, FALSE)
RETURNING
  *;

-- name: CreateOAuthUser :one
INSERT INTO
  users (username, email, name, google_id, is_active)
VALUES
  ($1, $2, $3, $4, TRUE)
RETURNING
  *;

-- name: UpdateUserProfile :execrows
UPDATE users
SET
  username = $2,
  email = $3
WHERE
  id = $1;

-- name: UpdateUserPassword :execrows
UPDATE users
SET
  password_hash = $2
WHERE
  id = $1;

-- name: UpdateUserIsActiveReturning :one
UPDATE users
SET
  is_active = $2
WHERE
  id = $1
RETURNING
  users.*;

-- name: UpdateUserGoogleID :execrows
UPDATE users
SET
  google_id = $2
WHERE
  id = $1
  AND google_id IS NULL;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE
  id = $1;
