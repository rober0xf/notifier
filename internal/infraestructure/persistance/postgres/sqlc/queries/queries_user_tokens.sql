-- name: CreateUserToken :one
INSERT INTO
  user_tokens (user_id, token_hash, purpose, expires_at)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: VerifyAndConsumeToken :one
UPDATE user_tokens
SET
  used = TRUE
WHERE
  token_hash = $1
  AND purpose = $2
  AND used = false
  AND expires_at > NOW()
RETURNING
  *;

-- name: GetTokenByHash :one
SELECT
  id,
  user_id,
  token_hash,
  purpose,
  used,
  expires_at,
  created_at
FROM
  user_tokens
WHERE
  token_hash = $1
  AND purpose = $2
  AND used = FALSE
  AND expires_at > NOW()
LIMIT
  1;

-- name: DeleteOldTokens :execrows
DELETE FROM user_tokens
WHERE
  used = TRUE
  OR expires_at < NOW();

-- name: DeleteByUserAndPurpose :exec
DELETE FROM user_tokens
WHERE
  user_id = $1
  AND purpose = $2;
