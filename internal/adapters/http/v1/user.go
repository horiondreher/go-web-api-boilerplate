package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

type SessionError struct {
	msg string
}

func (e *SessionError) Error() string {
	return e.msg
}

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

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) error {
	reqUser, err := httputils.Decode[CreateUserRequestDto](r)
	if err != nil {
		return errorResponse(err)
	}

	err = validate.Struct(reqUser)
	if err != nil {
		return errorResponse(err)
	}

	createdUser, err := adapter.userService.CreateUser(r.Context(), ports.NewUser{
		FullName: reqUser.FullName,
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})

	if err != nil {
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusCreated, CreateUserResponseDto{
		UID:      createdUser.UID,
		FullName: createdUser.FullName,
		Email:    createdUser.Email,
	})

	return nil
}

type LoginUserRequestDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponseDto struct {
	Email                 string    `json:"email"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (adapter *HTTPAdapter) loginUser(w http.ResponseWriter, r *http.Request) error {
	reqUser, err := httputils.Decode[LoginUserRequestDto](r)
	if err != nil {
		return errorResponse(err)
	}

	err = validate.Struct(reqUser)
	if err != nil {
		return errorResponse(err)
	}

	user, err := adapter.userService.LoginUser(r.Context(), ports.LoginUser{
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})
	if err != nil {
		return errorResponse(err)
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(user.Email, "user", adapter.config.AccessTokenDuration)
	if err != nil {
		return errorResponse(err)
	}

	refreshToken, refreshPayload, err := adapter.tokenMaker.CreateToken(user.Email, "user", adapter.config.RefreshTokenDuration)
	if err != nil {
		return errorResponse(err)
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
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusOK, loginRes)

	return nil
}

type RenewAccessTokenRequestDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RenewAccessTokenResponseDto struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (adapter *HTTPAdapter) renewAccessToken(w http.ResponseWriter, r *http.Request) error {
	renewAccessDto, err := httputils.Decode[RenewAccessTokenRequestDto](r)
	if err != nil {
		return errorResponse(err)
	}

	err = validate.Struct(renewAccessDto)
	if err != nil {
		return errorResponse(err)
	}

	refreshPayload, err := adapter.tokenMaker.VerifyToken(renewAccessDto.RefreshToken)
	if err != nil {
		return errorResponse(err)
	}

	session, err := adapter.userService.GetUserSession(r.Context(), refreshPayload.ID)
	if err != nil {
		return errorResponse(err)
	}

	if session.IsBlocked {
		return errorResponse(&SessionError{msg: "Session is blocked"})
	}

	if session.UserEmail != refreshPayload.Email {
		return errorResponse(&SessionError{msg: "Invalid session user"})
	}

	if session.RefreshToken != renewAccessDto.RefreshToken {
		return errorResponse(&SessionError{msg: "Invalid refresh token"})
	}

	if time.Now().After(session.ExpiresAt) {
		return errorResponse(&SessionError{msg: "Session expired"})
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(session.UserEmail, "user", adapter.config.AccessTokenDuration)

	if err != nil {
		return errorResponse(err)
	}

	renewAccessTokenResponse := RenewAccessTokenResponseDto{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	httputils.Encode(w, r, http.StatusOK, renewAccessTokenResponse)

	return nil
}

func (adapter *HTTPAdapter) getUserByUID(w http.ResponseWriter, r *http.Request) error {

	payload := r.Context().Value(middleware.KeyAuthUser).(*token.Payload)
	requestID := middleware.GetRequestID(r.Context())

	fmt.Println(payload)
	fmt.Println(requestID)

	userID := chi.URLParam(r, "uid")

	fmt.Println(userID)

	user, err := adapter.userService.GetUserByUID(r.Context(), userID)

	if err != nil {
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusOK, CreateUserResponseDto{
		UID:      user.UID,
		Email:    user.Email,
		FullName: user.FullName,
	})

	return nil
}
