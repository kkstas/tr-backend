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
	t.Parallel()

	t.Run("return 401 if no token", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)
		request := httptest.NewRequest("GET", "/user", nil)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns user data if user exists and token is valid", func(t *testing.T) {
		t.Parallel()
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
		testutils.AssertEqual(t, foundUser.ActiveVault, "")
	})

	t.Run("returns user data with activeVault if user has one", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)

		token, user := testutils.CreateUserWithToken(t, db)

		{
			vaultFC := struct {
				VaultName string `json:"vaultName"`
			}{VaultName: "asdf"}

			request := httptest.NewRequest("POST", "/vaults", testutils.ToJSONBuffer(t, vaultFC))
			request.Header.Set("Authorization", "Bearer "+token)
			response := httptest.NewRecorder()
			serv.ServeHTTP(response, request)
			testutils.AssertStatus(t, response.Code, http.StatusNoContent)
		}

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
		testutils.AssertNotEmpty(t, foundUser.ActiveVault)
	})
}
