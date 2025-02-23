package vault_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestDeleteOneByID(t *testing.T) {
	t.Parallel()

	t.Run("deletes existing vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		request := httptest.NewRequest("DELETE", "/vaults/"+vault.ID, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		vaults, err := testutils.NewTestVaultService(db).FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 0)
	})

	t.Run("does not delete existing vault if user is not assigned to the vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		_, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		otherUserToken, _ := testutils.CreateTestUserWithToken(t, db)
		request := httptest.NewRequest("DELETE", "/vaults/"+vault.ID, nil)
		request.Header.Set("Authorization", "Bearer "+otherUserToken)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)

		vaults, err := testutils.NewTestVaultService(db).FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 1)
	})
}
