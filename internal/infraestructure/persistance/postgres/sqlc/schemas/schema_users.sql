CREATE TYPE token_purpose AS ENUM ('email_verification', 'password_reset');

CREATE TABLE users (
    id serial PRIMARY KEY,
    username text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    PASSWORD text,
    name text,
    google_id text UNIQUE,
    active bool NOT NULL DEFAULT FALSE,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE user_tokens (
    id serial PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    purpose token_purpose NOT NULL,
    used bool NOT NULL DEFAULT FALSE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
