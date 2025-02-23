package expensecategory_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestCreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates expense category", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		serv, db := testutils.NewTestApplication(t)
		token, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		reqBody := struct {
			Name    string `json:"name"`
			VaultID string `json:"vaultID"`
		}{Name: "category name", VaultID: vault.ID}

		request := httptest.NewRequest("POST", "/expensecategories", testutils.ToJSONBuffer(t, reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, 204)

		foundCategories, err := testutils.NewTestExpenseCategoryService(db).FindAll(ctx, user.ID, reqBody.VaultID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(foundCategories), 1)
		testutils.AssertEqual(t, foundCategories[0].VaultID, reqBody.VaultID)
		testutils.AssertEqual(t, foundCategories[0].Name, reqBody.Name)
		testutils.AssertEqual(t, foundCategories[0].CreatedBy, user.ID)
	})

	t.Run("returns 400 if vault name is too short/long", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		reqBody := struct {
			Name    string `json:"name"`
			VaultID string `json:"vaultID"`
		}{Name: testutils.RandomString(1), VaultID: vault.ID}

		request := httptest.NewRequest("POST", "/expensecategories", testutils.ToJSONBuffer(t, reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

		var errorResBody = struct { // nolint: exhaustruct
			Name string `json:"name"`
		}{}
		err := json.NewDecoder(response.Body).Decode(&errorResBody)
		testutils.AssertNoError(t, err)
	})

	t.Run("returns 404 if user does not belong to given vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		_, _, vault := testutils.CreateTestUserWithTokenAndVault(t, db)
		token, _ := testutils.CreateTestUserWithToken(t, db)

		reqBody := struct {
			Name    string `json:"name"`
			VaultID string `json:"vaultID"`
		}{Name: testutils.RandomString(10), VaultID: vault.ID}

		request := httptest.NewRequest("POST", "/expensecategories", testutils.ToJSONBuffer(t, reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 404 if vault does not exist", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateTestUserWithToken(t, db)

		reqBody := struct {
			Name    string `json:"name"`
			VaultID string `json:"vaultID"`
		}{Name: testutils.RandomString(10), VaultID: uuid.New().String()}

		request := httptest.NewRequest("POST", "/expensecategories", testutils.ToJSONBuffer(t, reqBody))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
