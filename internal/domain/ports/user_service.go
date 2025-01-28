package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
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
	CreateUser(ctx context.Context, newUser NewUser) (pgsqlc.CreateUserRow, *domainerr.DomainError)
	LoginUser(ctx context.Context, loginUser LoginUser) (pgsqlc.User, *domainerr.DomainError)
	CreateUserSession(ctx context.Context, newUserSession NewUserSession) (pgsqlc.Session, *domainerr.DomainError)
	GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, *domainerr.DomainError)
	GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, *domainerr.DomainError)
}
