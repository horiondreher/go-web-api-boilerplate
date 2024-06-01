package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-boilerplate/internal/domain"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-boilerplate/internal/utils"
)

type UserManager struct {
	store pgsqlc.Querier
}

func NewUserManager(store pgsqlc.Querier) *UserManager {
	return &UserManager{
		store: store,
	}
}

func (service *UserManager) CreateUser(reqUser domain.CreateUserRequestDto) (domain.CreateUserResponseDto, error) {
	hashedPassword, err := utils.HashPassword(reqUser.Password)

	if err != nil {
		return domain.CreateUserResponseDto{}, err
	}

	args := pgsqlc.CreateUserParams{
		Email:     reqUser.Email,
		Password:  hashedPassword,
		FullName:  reqUser.FullName,
		IsStaff:   false,
		IsActive:  true,
		LastLogin: time.Now(),
	}

	ctx := context.Background()

	user, err := service.store.CreateUser(ctx, args)

	if err != nil {
		return domain.CreateUserResponseDto{}, err
	}

	return domain.CreateUserResponseDto{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (service *UserManager) LoginUser(reqUser domain.LoginUserRequestDto) (domain.LoginUserResponseDto, error) {
	user, err := service.store.GetUser(context.Background(), reqUser.Email)

	if err != nil {
		return domain.LoginUserResponseDto{}, err
	}

	err = utils.CheckPassword(reqUser.Password, user.Password)

	if err != nil {
		return domain.LoginUserResponseDto{}, err
	}

	return domain.LoginUserResponseDto{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (service *UserManager) CreateUserSession(refreshTokenID uuid.UUID, loggedUser *domain.LoginUserResponseDto, userAgent, clientIP string) (pgsqlc.Session, error) {
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

func (service *UserManager) GetUserSession(refreshTokenID uuid.UUID) (pgsqlc.Session, error) {
	session, err := service.store.GetSession(context.Background(), refreshTokenID)

	if err != nil {
		return pgsqlc.Session{}, err
	}

	return session, nil
}

func (service *UserManager) GetUserByUID(userUID string) (domain.LoginUserResponseDto, error) {

	parsedUID, err := uuid.Parse(userUID)

	if err != nil {
		return domain.LoginUserResponseDto{}, err
	}

	user, err := service.store.GetUserByUID(context.Background(), parsedUID)

	if err != nil {
		return domain.LoginUserResponseDto{}, err
	}

	return domain.LoginUserResponseDto{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
