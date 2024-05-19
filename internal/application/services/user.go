package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	store pgsqlc.Querier
}

func NewUserService(store pgsqlc.Querier) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) CreateUser(reqUser entities.CreateUserRequestDto) (entities.CreateUserResponseDto, error) {
	hashedPassword, err := utils.HashPassword(reqUser.Password)

	if err != nil {
		return entities.CreateUserResponseDto{}, err
	}

	args := pgsqlc.CreateUserParams{
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
		return entities.CreateUserResponseDto{}, err
	}

	return entities.CreateUserResponseDto{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (service *UserService) LoginUser(reqUser entities.LoginUserRequestDto) (entities.LoginUserResponseDto, error) {
	user, err := service.store.GetUser(context.Background(), reqUser.Email)

	if err != nil {
		return entities.LoginUserResponseDto{}, err
	}

	err = utils.CheckPassword(reqUser.Password, user.Password)

	if err != nil {
		return entities.LoginUserResponseDto{}, err
	}

	return entities.LoginUserResponseDto{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (service *UserService) CreateUserSession(refreshTokenID uuid.UUID, loggedUser *entities.LoginUserResponseDto, userAgent, clientIP string) (pgsqlc.Session, error) {
	session, err := service.store.CreateSession(context.Background(), pgsqlc.CreateSessionParams{
		Uid:          refreshTokenID,
		UserEmail:    loggedUser.Email,
		RefreshToken: loggedUser.RefreshToken,
		ExpiresAt:    loggedUser.RefreshTokenExpiresAt,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
	})

	if err != nil {
		return pgsqlc.Session{}, err
	}

	return session, nil
}

func (service *UserService) GetUserSession(refreshTokenID uuid.UUID) (pgsqlc.Session, error) {
	session, err := service.store.GetSession(context.Background(), refreshTokenID)

	if err != nil {
		return pgsqlc.Session{}, err
	}

	return session, nil
}
