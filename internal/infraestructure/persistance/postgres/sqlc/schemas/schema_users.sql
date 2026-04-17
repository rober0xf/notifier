CREATE TYPE token_purpose AS ENUM('email_verification', 'password_reset');

CREATE TYPE user_role AS ENUM('user', 'admin');

CREATE TABLE users (
  id serial PRIMARY KEY,
  username text NOT NULL UNIQUE,
  email text NOT NULL UNIQUE,
  password_hash text,
  name text,
  role user_role NOT NULL DEFAULT 'user',
  google_id text UNIQUE,
  is_active bool NOT NULL DEFAULT FALSE,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

ALTER TABLE users
ADD CONSTRAINT users_auth_method_check CHECK (
  password_hash IS NOT NULL
  OR google_id IS NOT NULL
);
