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

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) error {
	user, err := decode[entities.CreateUserRequestDto](r)
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

	encode(w, r, http.StatusCreated, userResponse)

	return nil
}

func (adapter *HTTPAdapter) loginUser(w http.ResponseWriter, r *http.Request) error {
	user, err := decode[entities.LoginUserRequestDto](r)
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

	encode(w, r, http.StatusOK, userResponse)

	return nil
}

func (adapter *HTTPAdapter) renewAccessToken(w http.ResponseWriter, r *http.Request) error {
	renewAccessDto, err := decode[entities.RenewAccessTokenRequestDto](r)
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

	renewAccessTokenResponse := entities.RenewAccessTokenResponseDto{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	encode(w, r, http.StatusOK, renewAccessTokenResponse)

	return nil
}
