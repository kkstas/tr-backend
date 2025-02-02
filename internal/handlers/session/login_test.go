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
}
