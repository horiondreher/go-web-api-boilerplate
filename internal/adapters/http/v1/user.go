package v1

import (
	"net/http"

	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
)

func (handler *HTTPHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := decode[entities.CreateUserRequestDto](r)

	if err != nil {
		error_response(w, r, err)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		error_response(w, r, err)
		return
	}

	userResponse, err := handler.service.CreateUser(user)

	if err != nil {
		error_response(w, r, err)
		return
	}

	encode(w, r, http.StatusCreated, userResponse)
}

func (handler *HTTPHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	user, err := decode[entities.LoginUserRequestDto](r)

	if err != nil {
		error_response(w, r, err)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		error_response(w, r, err)
		return
	}

	userResponse, err := handler.service.LoginUser(user)

	if err != nil {
		error_response(w, r, err)
		return
	}

	accessToken, accessPayload, err := handler.tokenMaker.CreateToken(userResponse.Email, "user", handler.config.AccessTokenDuration)

	if err != nil {
		error_response(w, r, err)
		return
	}

	refreshToken, refreshPayload, err := handler.tokenMaker.CreateToken(userResponse.Email, "user", handler.config.RefreshTokenDuration)
	if err != nil {
		error_response(w, r, err)
		return
	}

	userResponse.AccessToken = accessToken
	userResponse.AccessTokenExpiresAt = accessPayload.ExpiredAt
	userResponse.RefreshToken = refreshToken
	userResponse.RefreshTokenExpiresAt = refreshPayload.ExpiredAt

	encode(w, r, http.StatusOK, userResponse)
}
