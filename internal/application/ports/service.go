package ports

import (
	"github.com/google/uuid"
	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
)

type Service interface {
	CreateUser(user entities.CreateUserRequestDto) (entities.CreateUserResponseDto, error)
	LoginUser(user entities.LoginUserRequestDto) (entities.LoginUserResponseDto, error)
	CreateUserSession(refreshTokenID uuid.UUID, loggedUser *entities.LoginUserResponseDto, userAgent, clientIP string) (pgsqlc.Session, error)
	GetUserSession(refreshTokenID uuid.UUID) (pgsqlc.Session, error)
}
