package session_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestRegister(t *testing.T) {
	t.Run("registers new user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		userFC := struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}{
			Email:     "doe@johndoe.com",
			Password:  "mypassword123",
			FirstName: "John",
			LastName:  "Doe",
		}

		request := httptest.NewRequest("POST", "/register", testutils.ToJSONBuffer(t, userFC))
		response := httptest.NewRecorder()

		serv, cleanup, db := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		foundUsers, err := repositories.NewUserRepo(db).FindAll(ctx)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, foundUsers[0].Email, userFC.Email)
		testutils.AssertEqual(t, foundUsers[0].FirstName, userFC.FirstName)
		testutils.AssertEqual(t, foundUsers[0].LastName, userFC.LastName)
		testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
	})
}
