-- name: CreateSession :one
INSERT INTO "session" (
	"uid"
	, "user_email"
	, "refresh_token"
	, "user_agent"
	, "client_ip"
	, "is_blocked"
	, "expires_at"
	)
VALUES (
	$1
	, $2
	, $3
	, $4
	, $5
	, $6
	, $7
	) RETURNING *;

-- name: GetSession :one
SELECT *
FROM "session"
WHERE "uid" = $1 LIMIT 1;