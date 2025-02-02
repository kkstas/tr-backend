package session_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/auth"
	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestLogin(t *testing.T) {
	t.Run("returns auth token", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		password := testutils.RandomString(32)
		passwordHash, err := auth.HashPassword(password)
		testutils.AssertNoError(t, err)

		reqBody := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    "john@doe.eu",
			Password: password,
		}

		serv, cleanup, db := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)
		err = repositories.NewUserRepo(db).CreateOne(ctx, "John", "Doe", reqBody.Email, passwordHash)
		testutils.AssertNoError(t, err)

		request := httptest.NewRequest("POST", "/login", testutils.ToJSONBuffer(t, reqBody))
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var resBody auth.UserToken

		if err := json.NewDecoder(response.Body).Decode(&resBody); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		testutils.AssertEqual(t, resBody.TokenType, "Bearer")
		if len(resBody.Token) == 0 {
			t.Error("expected a token string in response, found empty string")
		}
	})

	t.Run("should reject invalid request properties", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, _ := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)

		type reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		tests := []struct {
			key    string
			reason string
			value  reqBody
		}{
			{key: "email", reason: "invalid email", value: reqBody{
				Email:    "doe@johndoe",
				Password: "mypassword123",
			}},
			{key: "password", reason: "too short", value: reqBody{
				Email:    "doe@johndoe.com",
				Password: "aa",
			}},
			{key: "password", reason: "too long", value: reqBody{
				Email:    "doe@johndoe.com",
				Password: testutils.RandomString(501),
			}},
		}

		for _, tc := range tests {
			t.Run(tc.key+tc.reason, func(t *testing.T) {
				t.Parallel()
				request := httptest.NewRequest("POST", "/login", testutils.ToJSONBuffer(t, tc.value))
				response := httptest.NewRecorder()
				serv.ServeHTTP(response, request)
				testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

				var m map[string]string
				if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				testutils.AssertNotEmpty(t, m[tc.key])
			})
		}
	})

	t.Run("returns 404 if user with given email does not exist", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, _ := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)

		reqBody := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{Email: "some@user.com", Password: "somepassword123"}

		request := httptest.NewRequest("POST", "/login", testutils.ToJSONBuffer(t, reqBody))
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
