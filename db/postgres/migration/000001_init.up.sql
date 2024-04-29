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


CREATE UNIQUE INDEX "user_email_idx" ON "user" USING BTREE ("email");

CREATE INDEX "user_uid_idx" ON "user" USING BTREE ("uid");
