package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain"
	"github.com/horiondreher/go-web-api-boilerplate/internal/infrastructure/persistence/pgsqlc"
)

type UserService interface {
	CreateUser(ctx context.Context, user domain.CreateUserRequestDto) (domain.CreateUserResponseDto, error)
	LoginUser(ctx context.Context, user domain.LoginUserRequestDto) (domain.LoginUserResponseDto, error)
	CreateUserSession(ctx context.Context, refreshTokenID uuid.UUID, loggedUser *domain.LoginUserResponseDto, userAgent, clientIP string) (pgsqlc.Session, error)
	GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, error)
	GetUserByUID(ctx context.Context, userUID string) (domain.LoginUserResponseDto, error)
}
