BEGIN;

-- Drop indexes
DROP INDEX "session_uid_idx";
DROP INDEX "user_uid_idx";
DROP INDEX "user_email_idx";

-- Drop constraints
ALTER TABLE "session" DROP CONSTRAINT "fk_user_id";

-- Drop tables
DROP TABLE "user";
DROP TABLE "session";

COMMIT;