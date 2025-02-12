package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestGetUserInfo(t *testing.T) {
	t.Run("return 401 if no token", func(t *testing.T) {
		serv, _ := testutils.NewTestApplication(t)
		request := httptest.NewRequest("GET", "/user", nil)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns user data if user exists and token is valid", func(t *testing.T) {
		serv, db := testutils.NewTestApplication(t)

		token, user := testutils.CreateUserWithToken(t, db)

		request := httptest.NewRequest("GET", "/user", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var foundUser models.User
		if err := json.NewDecoder(response.Body).Decode(&foundUser); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		testutils.AssertEqual(t, foundUser.ID, user.ID)
		testutils.AssertEqual(t, foundUser.FirstName, user.FirstName)
		testutils.AssertEqual(t, foundUser.LastName, user.LastName)
		testutils.AssertEqual(t, foundUser.Email, user.Email)
		testutils.AssertEqual(t, foundUser.CreatedAt, user.CreatedAt)
	})
}
