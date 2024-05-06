package v1

import (
	"net/http"
	"time"

	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
)

type SessionError struct {
	msg string
}

func (e *SessionError) Error() string {
	return e.msg
}

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := decode[entities.CreateUserRequestDto](r)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	userResponse, err := adapter.userService.CreateUser(user)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	encode(w, r, http.StatusCreated, userResponse)
}

func (adapter *HTTPAdapter) loginUser(w http.ResponseWriter, r *http.Request) {
	user, err := decode[entities.LoginUserRequestDto](r)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	userResponse, err := adapter.userService.LoginUser(user)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(userResponse.Email, "user", adapter.config.AccessTokenDuration)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	refreshToken, refreshPayload, err := adapter.tokenMaker.CreateToken(userResponse.Email, "user", adapter.config.RefreshTokenDuration)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	userResponse.AccessToken = accessToken
	userResponse.AccessTokenExpiresAt = accessPayload.ExpiredAt
	userResponse.RefreshToken = refreshToken
	userResponse.RefreshTokenExpiresAt = refreshPayload.ExpiredAt

	_, err = adapter.userService.CreateUserSession(refreshPayload.ID, &userResponse, r.UserAgent(), r.RemoteAddr)

	if err != nil {
		errorResponse(w, r, err)
		return
	}

	encode(w, r, http.StatusOK, userResponse)
}

func (adapter *HTTPAdapter) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	renewAccessDto, err := decode[entities.RenewAccessTokenRequestDto](r)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	err = validate.Struct(renewAccessDto)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	refreshPayload, err := adapter.tokenMaker.VerifyToken(renewAccessDto.RefreshToken)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	session, err := adapter.userService.GetUserSession(refreshPayload.ID)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	if session.IsBlocked {
		errorResponse(w, r, &SessionError{msg: "Session is blocked"})
		return
	}

	if session.UserEmail != refreshPayload.Email {
		errorResponse(w, r, &SessionError{msg: "Invalid session user"})
		return
	}

	if session.RefreshToken != renewAccessDto.RefreshToken {
		errorResponse(w, r, &SessionError{msg: "Invalid refresh token"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		errorResponse(w, r, &SessionError{msg: "Session expired"})
		return
	}

	accessToken, accessPayload, err := adapter.tokenMaker.CreateToken(session.UserEmail, "user", adapter.config.AccessTokenDuration)

	if err != nil {
		errorResponse(w, r, err)
		return
	}

	renewAccessTokenResponse := entities.RenewAccessTokenResponseDto{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	encode(w, r, http.StatusOK, renewAccessTokenResponse)
}
