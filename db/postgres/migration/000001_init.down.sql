BEGIN;

-- Drop constraints
ALTER TABLE "session" DROP CONSTRAINT "fk_user_email";

-- Drop indexes
DROP INDEX "session_uid_idx";
DROP INDEX "user_uid_idx";
DROP INDEX "user_email_idx";

-- Drop tables
DROP TABLE "session";
DROP TABLE "user";

COMMIT;