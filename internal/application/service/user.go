package service

import (
	"context"
	"time"

	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/postgres"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store postgres.Querier
}

func NewService(store postgres.Querier) *Service {
	return &Service{
		store: store,
	}
}

func (service *Service) CreateUser(reqUser entities.CreateUserRequestDto) (entities.CreateUserResponseDto, error) {
	hashedPassword, err := utils.HashPassword(reqUser.Password)

	if err != nil {
		log.Err(err).Msg("Error hashing password")
		return entities.CreateUserResponseDto{}, err
	}

	args := postgres.CreateUserParams{
		Email:    reqUser.Email,
		Password: hashedPassword,
		FullName: reqUser.FullName,
		IsStaff:  false,
		IsActive: true,
		LastLogin: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	ctx := context.Background()

	user, err := service.store.CreateUser(ctx, args)

	if err != nil {
		log.Err(err).Msg("Error creating user")
		return entities.CreateUserResponseDto{}, err
	}

	return entities.CreateUserResponseDto{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (service *Service) LoginUser(reqUser entities.LoginUserRequestDto) (entities.LoginUserResponseDto, error) {
	user, err := service.store.GetUser(context.Background(), reqUser.Email)

	if err != nil {
		log.Err(err).Msg("Error getting user by email")
		return entities.LoginUserResponseDto{}, err
	}

	err = utils.CheckPassword(reqUser.Password, user.Password)

	if err != nil {
		log.Err(err).Msg("Invalid password")
		return entities.LoginUserResponseDto{}, err
	}

	return entities.LoginUserResponseDto{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
