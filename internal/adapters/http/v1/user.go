package v1

import (
	"net/http"

	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
)

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

	userResponse, err := adapter.service.CreateUser(user)

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

	userResponse, err := adapter.service.LoginUser(user)

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

	encode(w, r, http.StatusOK, userResponse)
}
