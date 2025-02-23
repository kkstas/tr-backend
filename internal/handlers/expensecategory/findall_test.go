package expensecategory_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestFindAll(t *testing.T) {
	t.Parallel()

	t.Run("finds all expense categories", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		serv, db := testutils.NewTestApplication(t)
		token, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		err := testutils.NewTestExpenseCategoryService(db).CreateOne(ctx, "some name", user.ID, vault.ID)
		testutils.AssertNoError(t, err)

		request := httptest.NewRequest("GET", "/expensecategories/"+vault.ID, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, 200)

		var categories []models.ExpenseCategory

		err = json.NewDecoder(response.Body).Decode(&categories)
		testutils.AssertNoError(t, err)
	})

	t.Run("returns 404 if there's no vault with provided id", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateTestUserWithToken(t, db)

		request := httptest.NewRequest("GET", "/expensecategories/"+uuid.New().String(), nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 404 if user does not belong to vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		token, _ := testutils.CreateTestUserWithToken(t, db)
		_, _, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		request := httptest.NewRequest("GET", "/expensecategories/"+vault.ID, nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
