package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestCreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates new vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, user := testutils.CreateUserWithToken(t, db)

		vaultFC := struct {
			VaultName string `json:"vaultName"`
		}{VaultName: "asdf"}

		request := httptest.NewRequest("POST", "/vaults", testutils.ToJSONBuffer(t, vaultFC))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		vaults, err := testutils.NewTestVaultService(db).FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(vaults), 1)
		testutils.AssertEqual(t, vaults[0].Name, vaultFC.VaultName)
		testutils.AssertEqual(t, vaults[0].UserRole, models.VaultRoleOwner)
	})

	t.Run("vault id is saved as user's active vault if user had no active vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, user := testutils.CreateUserWithToken(t, db)

		vaultFC := struct {
			VaultName string `json:"vaultName"`
		}{VaultName: "asdf"}

		request := httptest.NewRequest("POST", "/vaults", testutils.ToJSONBuffer(t, vaultFC))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		vaults, err := testutils.NewTestVaultService(db).FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(vaults), 1)
		testutils.AssertEqual(t, vaults[0].Name, vaultFC.VaultName)
		testutils.AssertEqual(t, vaults[0].UserRole, models.VaultRoleOwner)

		newUserData, err := testutils.NewTestUserService(db).FindOneByID(context.Background(), user.ID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, newUserData.ActiveVault, vaults[0].ID)
	})

	t.Run("returns 400 with error message if request body is invalid", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateUserWithToken(t, db)
		request := httptest.NewRequest("POST", "/vaults",
			testutils.ToJSONBuffer(t, struct {
				VaultName string `json:"vaultName"`
			}{VaultName: ""}),
		)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

		type resBody struct {
			VaultName string `json:"vaultName"`
		}
		var body resBody

		err := json.NewDecoder(response.Body).Decode(&body)
		testutils.AssertNoError(t, err)

		if len(body.VaultName) == 0 {
			t.Error("expected vaultName error message in response body but didn't get one")
		}
	})

	t.Run("returns 400 with error message in response body if no body in the request", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateUserWithToken(t, db)
		request := httptest.NewRequest("POST", "/vaults", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

		type resBody struct {
			Message string `json:"message"`
		}
		var body resBody

		err := json.NewDecoder(response.Body).Decode(&body)
		testutils.AssertNoError(t, err)

		if len(body.Message) == 0 {
			t.Error("expected error message in response body but didn't get one")
		}
	})

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)
		request := httptest.NewRequest("POST", "/vaults", nil)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})
}
