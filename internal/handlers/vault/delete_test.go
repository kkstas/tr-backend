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
		token, user := testutils.CreateTestUserWithToken(t, db)

		vaultService := testutils.NewTestVaultService(db)
		err := vaultService.CreateOne(context.Background(), user.ID, "some vault name")
		testutils.AssertNoError(t, err)

		vaults, err := vaultService.FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 1)

		request := httptest.NewRequest("DELETE", "/vaults/"+vaults[0].ID, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		vaults, err = vaultService.FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 0)
	})

	t.Run("does not delete existing vault if user is not assigned to the vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		user := testutils.CreateTestUser(t, db)

		vaultService := testutils.NewTestVaultService(db)
		err := vaultService.CreateOne(context.Background(), user.ID, "some vault name")
		testutils.AssertNoError(t, err)

		vaults, err := vaultService.FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 1)

		otherUserToken, _ := testutils.CreateTestUserWithToken(t, db)
		request := httptest.NewRequest("DELETE", "/vaults/"+vaults[0].ID, nil)
		request.Header.Set("Authorization", "Bearer "+otherUserToken)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)

		vaults, err = vaultService.FindAll(context.Background(), user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(vaults), 1)
	})
}
