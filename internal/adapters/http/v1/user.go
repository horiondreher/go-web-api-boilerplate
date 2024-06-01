package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/horiondreher/go-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-boilerplate/internal/domain"
)

type SessionError struct {
	msg string
}

func (e *SessionError) Error() string {
	return e.msg
}

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) error {
	user, err := httputils.Decode[domain.CreateUserRequestDto](r)
	if err != nil {
		return errorResponse(err)
	}

	err = validate.Struct(user)
	if err != nil {
		return errorResponse(err)
	}

	userResponse, err := adapter.userService.CreateUser(user)
	if err != nil {
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusCreated, userResponse)

	return nil
}

func (adapter *HTTPAdapter) loginUser(w http.ResponseWriter, r *http.Request) error {
	user, err := httputils.Decode[domain.LoginUserRequestDto](r)
	if err != nil {
		return errorResponse(err)
	}

	err = validate.Struct(user)
	if err != nil {
		return errorResponse(err)
	}

	userResponse, err := adapter.userService.LoginUser(user)
	if err != nil {
		return errorResponse(err)
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(userResponse.Email, "user", adapter.config.AccessTokenDuration)
	if err != nil {
		return errorResponse(err)
	}

	refreshToken, refreshPayload, err := adapter.tokenMaker.CreateToken(userResponse.Email, "user", adapter.config.RefreshTokenDuration)
	if err != nil {
		return errorResponse(err)
	}

	userResponse.AccessToken = accessToken
	userResponse.AccessTokenExpiresAt = accessPayload.ExpiredAt
	userResponse.RefreshToken = refreshToken
	userResponse.RefreshTokenExpiresAt = refreshPayload.ExpiredAt

	_, err = adapter.userService.CreateUserSession(refreshPayload.ID, &userResponse, r.UserAgent(), r.RemoteAddr)

	if err != nil {
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusOK, userResponse)

	return nil
}

func (adapter *HTTPAdapter) renewAccessToken(w http.ResponseWriter, r *http.Request) error {
	renewAccessDto, err := httputils.Decode[domain.RenewAccessTokenRequestDto](r)
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

	session, err := adapter.userService.GetUserSession(refreshPayload.ID)
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

	renewAccessTokenResponse := domain.RenewAccessTokenResponseDto{
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

	user, err := adapter.userService.GetUserByUID(userID)

	if err != nil {
		return errorResponse(err)
	}

	httputils.Encode(w, r, http.StatusOK, user)

	return nil
}
