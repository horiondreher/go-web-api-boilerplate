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
			name: "Create user",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s", "password": "%s"}`, user.full_name, user.email, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				validateUserResponse(t, user, recorder.Body)
			},
		},
		{
			name: "Create user with invalid email",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s", "password": "%s"}`, user.full_name, "invalid_email", user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create user without name",
			body: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.email, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create user without email",
			body: fmt.Sprintf(`{"full_name": "%s", "password": "%s"}`, user.full_name, user.password),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create user without password",
			body: fmt.Sprintf(`{"full_name": "%s", "email": "%s"}`, user.full_name, user.email),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create user with empty body",
			body: `{}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create user with invalid json",
			body: `{"full_name": "invalid_json}`,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server, err := NewHTTPHandler(testService)

			require.NoError(t, err)

			server.createUser(recorder, req)

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
