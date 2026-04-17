CREATE TABLE user_tokens (
  id serial PRIMARY KEY,
  user_id int NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  token_hash text NOT NULL,
  purpose token_purpose NOT NULL,
  used bool NOT NULL DEFAULT FALSE,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW()
);
