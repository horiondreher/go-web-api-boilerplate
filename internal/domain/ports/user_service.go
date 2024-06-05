package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/infrastructure/persistence/pgsqlc"
)

type NewUser struct {
	FullName string
	Email    string
	Password string
}

type LoginUser struct {
	Email    string
	Password string
}

type NewUserSession struct {
	RefreshTokenID        uuid.UUID
	Email                 string
	RefreshToken          string
	UserAgent             string
	ClientIP              string
	RefreshTokenExpiresAt time.Time
}

type UserService interface {
	CreateUser(ctx context.Context, newUser NewUser) (pgsqlc.CreateUserRow, error)
	LoginUser(ctx context.Context, loginUser LoginUser) (pgsqlc.User, error)
	CreateUserSession(ctx context.Context, newUserSession NewUserSession) (pgsqlc.Session, error)
	GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, error)
	GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, error)
}
