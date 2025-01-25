package services

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
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

func (service *UserManager) CreateUser(ctx context.Context, newUser ports.NewUser) (pgsqlc.CreateUserRow, *domainerr.DomainError) {
	hashedPassword, hashErr := utils.HashPassword(newUser.Password)
	if hashErr != nil {
		return pgsqlc.CreateUserRow{}, hashErr
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
	if err != nil {
		return pgsqlc.CreateUserRow{}, domainerr.MatchPostgresError(err)
	}

	return user, nil
}

func (service *UserManager) LoginUser(ctx context.Context, loginUser ports.LoginUser) (pgsqlc.User, *domainerr.DomainError) {
	user, err := service.store.GetUser(ctx, loginUser.Email)
	if err != nil {
		return pgsqlc.User{}, domainerr.MatchPostgresError(err)
	}

	passErr := utils.CheckPassword(loginUser.Password, user.Password)
	if passErr != nil {
		return pgsqlc.User{}, passErr
	}

	return user, nil
}

func (service *UserManager) CreateUserSession(ctx context.Context, newUserSession ports.NewUserSession) (pgsqlc.Session, *domainerr.DomainError) {
	session, err := service.store.CreateSession(ctx, pgsqlc.CreateSessionParams{
		UID:          newUserSession.RefreshTokenID,
		UserEmail:    newUserSession.Email,
		RefreshToken: newUserSession.RefreshToken,
		ExpiresAt:    newUserSession.RefreshTokenExpiresAt,
		UserAgent:    newUserSession.UserAgent,
		ClientIP:     newUserSession.ClientIP,
	})
	if err != nil {
		return pgsqlc.Session{}, domainerr.MatchPostgresError(err)
	}

	return session, nil
}

func (service *UserManager) GetUserSession(ctx context.Context, refreshTokenID uuid.UUID) (pgsqlc.Session, *domainerr.DomainError) {
	session, err := service.store.GetSession(ctx, refreshTokenID)
	if err != nil {
		return pgsqlc.Session{}, domainerr.MatchPostgresError(err)
	}

	return session, nil
}

func (service *UserManager) GetUserByUID(ctx context.Context, userUID string) (pgsqlc.User, *domainerr.DomainError) {
	parsedUID, err := uuid.Parse(userUID)
	if err != nil {
		return pgsqlc.User{}, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.UnexpectedError, err.Error(), err)
	}

	user, err := service.store.GetUserByUID(ctx, parsedUID)
	if err != nil {
		return pgsqlc.User{}, domainerr.MatchPostgresError(err)
	}

	return user, nil
}
