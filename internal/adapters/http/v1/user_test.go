package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horiondreher/go-boilerplate/internal/domain/entities"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/stretchr/testify/require"
)

type testUser struct {
	full_name string
	email     string
	password  string
}

func TestCreateUserV1(t *testing.T) {
	user := testUser{
		full_name: utils.RandomString(6),
		email:     utils.RandomEmail(),
		password:  utils.RandomString(6),
	}

	tt := []struct {
		name          string
		body          string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "CreateUser",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s", "password": "%s"}`, user.full_name, user.email, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				validateUserResponse(t, user, recorder.Body)
			},
		},
		{
			name: "CreateUserWithInvalidEmail",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s", "password": "%s"}`, user.full_name, "invalid_email", user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "CreateUserWithoutName",
			body: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.email, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "CreateUserWithoutEmail",
			body: fmt.Sprintf(`{"full_name": "%s", "password": "%s"}`, user.full_name, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "CreateUserWithoutPassword",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s"}`, user.full_name, user.email),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "CreateUserWithEmptyBody",
			body: `{}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "CreateUserWithInvalidJson",
			body: `{"full_name": "invalid_json}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "CreateUserWithEmptyJson",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server, err := NewHTTPAdapter(testUserService)

			require.NoError(t, err)

			handlerFn := server.handlerWrapper(server.createUser)
			handlerFn(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func validateUserResponse(t *testing.T, response testUser, body *bytes.Buffer) {
	var responseUser entities.CreateUserResponseDto
	err := json.NewDecoder(body).Decode(&responseUser)
	require.NoError(t, err)

	require.Equal(t, response.full_name, responseUser.FullName)
	require.Equal(t, response.email, responseUser.Email)

	require.NotZero(t, responseUser.ID)
	require.IsType(t, int64(0), responseUser.ID)
}

func TestLoginUser(t *testing.T) {
	user := testUser{
		full_name: utils.RandomString(6),
		email:     utils.RandomEmail(),
		password:  utils.RandomString(6),
	}

	_, err := testUserService.CreateUser(entities.CreateUserRequestDto{
		FullName: user.full_name,
		Email:    user.email,
		Password: user.password,
	})

	require.NoError(t, err)

	tt := []struct {
		name          string
		body          string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "LoginUser",
			body: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.email, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "LoginUserWithInvalidEmail",
			body: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "invalid_email", user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "LoginUserWithoutEmail",
			body: fmt.Sprintf(`{"password": "%s"}`, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "LoginUserWithEmptyPassword",
			body: fmt.Sprintf(`{"email": "%s"}`, user.email),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "LoginUserWithEmptyBody",
			body: `{}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "LoginUserWithInvalidJson",
			body: `{"email": "invalid_json}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "LoginUserWithEmptyJson",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server, err := NewHTTPAdapter(testUserService)

			require.NoError(t, err)

			handlerFn := server.handlerWrapper(server.loginUser)
			handlerFn(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}
