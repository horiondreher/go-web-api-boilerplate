// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: session.sql

package pgsqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
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
	) RETURNING id, uid, user_email, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
`

type CreateSessionParams struct {
	Uid          uuid.UUID
	UserEmail    string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.Uid,
		arg.UserEmail,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.UserEmail,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, uid, user_email, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
FROM "session"
WHERE "uid" = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, uid uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, uid)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.UserEmail,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
