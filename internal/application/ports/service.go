package ports

import "github.com/horiondreher/go-boilerplate/internal/domain/entities"

type Service interface {
	CreateUser(user entities.CreateUserRequestDto) (entities.CreateUserResponseDto, error)
	LoginUser(user entities.LoginUserRequestDto) (entities.LoginUserResponseDto, error)
}
