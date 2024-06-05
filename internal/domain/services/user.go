package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/horiondreher/go-web-api-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"
)

type UserManager struct {
	store pgsqlc.Querier
}

func NewUserManager(store pgsqlc.Querier) *UserManager {
	return &UserManager{
		store: store,
	}
}

func (service *UserManager) CreateUser(ctx context.Context, newUser ports.NewUser) (pgsqlc.CreateUserRow, error) {
	hashedPassword, err := utils.HashPassword(newUser.Password)

	if err != nil {
		return pgsqlc.CreateUserRow{}, err
	}

	args := pgsqlc.CreateUserParams{
		Email:     newUser.Email,
		Password:  hashedPassword,
		FullName:  newUser.FullName,
		IsStaff:   false,
		IsActive:  true,
		LastLogin: time.Now(),
	}

	user, err := service.store.CreateUser(ctx, args)

	return user, err
}

func (service *UserManager) LoginUser(ctx context.Context, loginUser ports.LoginUser) (pgsqlc.User, error) {
	user, err := service.store.GetUser(ctx, loginUser.Email)

	if err != nil {
		return pgsqlc.User{}, err
	}

	err = utils.CheckPassword(loginUser.Password, user.Password)

	if err != nil {
		return pgsqlc.User{}, err
	}

	return user, nil
}

func (service *UserManager) CreateUserSession(ctx context.Context, newUserSession ports.NewUserSession) (pgsqlc.Session, error) {
	session, err := service.store.CreateSession(ctx, pgsqlc.CreateSessionParams{
		UID:          newUserSession.RefreshTokenID,
		UserEmail:    newUserSession.Email,
		RefreshToken: newUserSession.RefreshToken,
		ExpiresAt:    newUserSession.RefreshTokenExpiresAt,
		UserAgent:    newUserSession.UserAgent,
		ClientIP:     newUserSession.ClientIP,
	})

	if err != nil {
		return pgsqlc.Session{}, err
	}

	return session, nil
}

func (service *UserManager) GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, error) {
	session, err := service.store.GetSession(ctx, refreshTokenID)

	if err != nil {
		return pgsqlc.Session{}, err
	}

	return session, nil
}

func (service *UserManager) GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, error) {

	parsedUID, err := uuid.Parse(userUID)

	if err != nil {
		return pgsqlc.User{}, err
	}

	user, err := service.store.GetUserByUID(ctx, parsedUID)

	return user, err
}
