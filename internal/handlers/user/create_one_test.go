package user_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateOneUser(t *testing.T) {
	t.Run("returns status 204 & saves user in DB when created new user", func(t *testing.T) {
		serv, cleanup, _ := testutils.NewTestApplication(t)
		defer cleanup()

		reqBody := testutils.ToJSONBuffer(t, models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@email.com",
		})

		response := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/users", reqBody)
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		// check if users have been saved
		response = httptest.NewRecorder()
		serv.ServeHTTP(response, httptest.NewRequest("GET", "/users", nil))
		foundUsers := testutils.DecodeJSON[[]models.User](t, response.Body)
		want := 1
		if len(foundUsers) != want {
			t.Errorf("got %d users, want %d", len(foundUsers), want)
		}
	})
}
