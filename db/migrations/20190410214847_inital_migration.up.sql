CREATE EXTENSION IF NOT EXISTS "pgcrypto";

create table users (
  uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text UNIQUE NOT NULL,
  name text NOT NULL,
  password text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz
);

create table sessions (
  user_uuid uuid NOT NULL REFERENCES users (uuid),
  created_at timestamptz NOT NULL,
  expires_at timestamptz NOT NULL,
  deleted_at timestamptz NOT NULL
);