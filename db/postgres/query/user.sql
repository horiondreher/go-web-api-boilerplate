-- name: CreateUser :one
INSERT INTO "user" (
        "email",
        "password",
        "full_name",
        "is_staff",
        "is_active",
        "last_login"
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING "id", "email", "full_name", "created_at", "modified_at";

-- name: GetUser :one
SELECT *
FROM "user"
WHERE "email" = $1
LIMIT 1;