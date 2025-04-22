package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/rs/zerolog/log"
)

type CreateUserRequestDto struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponseDto struct {
	UID      uuid.UUID `json:"uid"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	reqUser, err := httputils.Decode[CreateUserRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqUser)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	createdUser, err := adapter.userService.CreateUser(r.Context(), ports.NewUser{
		FullName: reqUser.FullName,
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})
	if err != nil {
		return err
	}

	log.Info().Msg("AAAAAAAAAAAA")

	err = httputils.Encode(w, r, http.StatusCreated, CreateUserResponseDto{
		UID:      createdUser.UID,
		FullName: createdUser.FullName,
		Email:    createdUser.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

type LoginUserRequestDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponseDto struct {
	Email                 string    `json:"email"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (adapter *HTTPAdapter) loginUser(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	reqUser, err := httputils.Decode[LoginUserRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqUser)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	user, err := adapter.userService.LoginUser(r.Context(), ports.LoginUser{
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})
	if err != nil {
		return err
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(user.Email, "user", adapter.config.AccessTokenDuration)
	if err != nil {
		return err
	}

	refreshToken, refreshPayload, err := adapter.tokenMaker.CreateToken(user.Email, "user", adapter.config.RefreshTokenDuration)
	if err != nil {
		return err
	}

	loginRes := LoginUserResponseDto{
		Email:                 user.Email,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	_, err = adapter.userService.CreateUserSession(r.Context(), ports.NewUserSession{
		RefreshTokenID:        refreshPayload.ID,
		Email:                 loginRes.Email,
		RefreshToken:          loginRes.RefreshToken,
		RefreshTokenExpiresAt: loginRes.RefreshTokenExpiresAt,
		UserAgent:             r.UserAgent(),
		ClientIP:              r.RemoteAddr,
	})
	if err != nil {
		return err
	}

	err = httputils.Encode(w, r, http.StatusOK, loginRes)
	if err != nil {
		return err
	}

	return nil
}

type RenewAccessTokenRequestDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RenewAccessTokenResponseDto struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (adapter *HTTPAdapter) renewAccessToken(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	renewAccessDto, err := httputils.Decode[RenewAccessTokenRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(renewAccessDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	refreshPayload, err := adapter.tokenMaker.VerifyToken(renewAccessDto.RefreshToken)
	if err != nil {
		return err
	}

	session, err := adapter.userService.GetUserSession(r.Context(), refreshPayload.ID)
	if err != nil {
		return err
	}

	if session.IsBlocked {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "session is blocked", errors.New("session is blocked"))
	}

	if session.UserEmail != refreshPayload.Email {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid session user", errors.New("invalid session user"))
	}

	if session.RefreshToken != renewAccessDto.RefreshToken {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid refresh token", errors.New("invalid refresh token"))
	}

	if time.Now().After(session.ExpiresAt) {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "session expired", errors.New("session expired"))
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(session.UserEmail, "user", adapter.config.AccessTokenDuration)
	if err != nil {
		return err
	}

	renewAccessTokenResponse := RenewAccessTokenResponseDto{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	err = httputils.Encode(w, r, http.StatusOK, renewAccessTokenResponse)
	if err != nil {
		return err
	}

	return nil
}

func (adapter *HTTPAdapter) getUserByUID(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	payload := r.Context().Value(middleware.KeyAuthUser).(*token.Payload)
	requestID := middleware.GetRequestID(r.Context())

	fmt.Println(payload)
	fmt.Println(requestID)

	userID := chi.URLParam(r, "uid")

	user, serviceErr := adapter.userService.GetUserByUID(r.Context(), userID)
	if serviceErr != nil {
		return serviceErr
	}

	err := httputils.Encode(w, r, http.StatusOK, CreateUserResponseDto{
		UID:      user.UID,
		Email:    user.Email,
		FullName: user.FullName,
	})
	if err != nil {
		return err
	}

	return nil
}
