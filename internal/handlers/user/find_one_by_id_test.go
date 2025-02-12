package user_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestFindOneUser(t *testing.T) {
	t.Run("returns status 200 & found user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		serv, db := testutils.NewTestApplication(t)

		userRepo := repositories.NewUserRepo(db)

		firstName := "John"
		lastName := "Doe"
		email := "john@doe.com"

		err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
		if err != nil {
			t.Fatalf("failed to create new user in repo: %v", err)
		}

		foundUsers, err := userRepo.FindAll(ctx)
		if err != nil {
			t.Fatalf("failed to find users in repo: %v", err)
		}

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/users/"+foundUsers[0].ID, nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var foundUser models.User
		if err := json.NewDecoder(response.Body).Decode(&foundUser); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		testutils.AssertEqual(t, foundUser.FirstName, firstName)
		testutils.AssertEqual(t, foundUser.LastName, lastName)
		testutils.AssertEqual(t, foundUser.Email, email)
		testutils.AssertValidDate(t, foundUser.CreatedAt)
		if err := uuid.Validate(foundUser.ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})

	t.Run("returns status 400 if id is not valid uuid", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/users/asdfasdf", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns status 404 if no user is found", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/users/ff16bbc6-d671-4fc3-9ba6-7073122f715c", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
