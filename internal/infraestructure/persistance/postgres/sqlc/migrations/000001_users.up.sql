CREATE TABLE users (
    id serial PRIMARY KEY,
    username text NOT NULL,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    active bool NOT NULL DEFAULT FALSE,
    email_verification_hash text,
    created_at timestamptz NOT NULL DEFAULT now(),
    timeout interval
)
