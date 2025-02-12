package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestFindAll(t *testing.T) {
	t.Run("finds all vaults that belong to user making the request", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)

		// create vault for user
		token, user := testutils.CreateUserWithToken(t, db)
		err := repositories.NewVaultRepo(db).CreateOne(context.Background(), user.ID, "name")
		testutils.AssertNoError(t, err)

		// also create vault for other user that should not be returned
		{
			_, otherUser := testutils.CreateUserWithToken(t, db)
			err := repositories.NewVaultRepo(db).CreateOne(context.Background(), otherUser.ID, "asdf")
			testutils.AssertNoError(t, err)
		}

		request := httptest.NewRequest("GET", "/vaults", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var vaults = []models.UserVaultWithRole{}
		err = json.NewDecoder(response.Body).Decode(&vaults)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(vaults), 1)
	})

	t.Run("returns empty array if no user vaults are found", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateUserWithToken(t, db)

		request := httptest.NewRequest("GET", "/vaults", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var vaults = []models.UserVaultWithRole{}
		err := json.NewDecoder(response.Body).Decode(&vaults)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(vaults), 0)
	})

	t.Run("returns 401 if unauthorized", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)
		request := httptest.NewRequest("GET", "/vaults", nil)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})
}
