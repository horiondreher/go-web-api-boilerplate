BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tables
CREATE TABLE "user" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "uid" UUID DEFAULT uuid_generate_v4(),
  "email" VARCHAR(255) NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "full_name" VARCHAR(255) NOT NULL,
  "is_staff" BOOLEAN NOT NULL,
  "is_active" BOOLEAN NOT NULL,
  "last_login" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "modified_at" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "session" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "uid" UUID NOT NULL,
  "user_id" BIGINT NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT false,
  "expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

-- Add constraints
ALTER TABLE "session"
ADD CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE;

-- Add indexes
CREATE UNIQUE INDEX "user_email_idx" ON "user" USING BTREE ("email");
CREATE INDEX "user_uid_idx" ON "user" USING BTREE ("uid");
CREATE INDEX "session_uid_idx" ON "session" USING BTREE ("uid");

COMMIT;