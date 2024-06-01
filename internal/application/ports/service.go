package ports

import (
	"github.com/google/uuid"
	"github.com/horiondreher/go-boilerplate/internal/domain"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
)

type UserService interface {
	CreateUser(user domain.CreateUserRequestDto) (domain.CreateUserResponseDto, error)
	LoginUser(user domain.LoginUserRequestDto) (domain.LoginUserResponseDto, error)
	CreateUserSession(refreshTokenID uuid.UUID, loggedUser *domain.LoginUserResponseDto, userAgent, clientIP string) (pgsqlc.Session, error)
	GetUserSession(refreshTokenID uuid.UUID) (pgsqlc.Session, error)
	GetUserByUID(userUID string) (domain.LoginUserResponseDto, error)
}
