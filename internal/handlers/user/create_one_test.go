package user_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateOneUser(t *testing.T) {
	t.Run("returns status 204 & saves user in DB when created new user", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, _ := testutils.NewTestApplication(t)
		defer cleanup()

		userFC := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@email.com",
		}

		reqBody := testutils.ToJSONBuffer(t, userFC)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/users", reqBody)
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		// Check if user data has been saved properly
		response = httptest.NewRecorder()
		serv.ServeHTTP(response, httptest.NewRequest("GET", "/users", nil))
		foundUsers := testutils.DecodeJSON[[]models.User](t, response.Body)
		want := 1
		if len(foundUsers) != want {
			t.Errorf("got %d users, want %d", len(foundUsers), want)
		}

		testutils.AssertEqual(t, foundUsers[0].FirstName, userFC.FirstName)
		testutils.AssertEqual(t, foundUsers[0].LastName, userFC.LastName)
		testutils.AssertEqual(t, foundUsers[0].Email, userFC.Email)
		testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
		if err := uuid.Validate(foundUsers[0].ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}

	})
}
