package user_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestFindAllUsers(t *testing.T) {
	t.Run("returns status 200 & array with users", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, db := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)

		err := repositories.NewUserRepo(db).CreateOne(context.Background(), "John", "Doe", "john@doe.com")
		if err != nil {
			t.Fatalf("failed to create new user in repo: %v", err)
		}

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/users", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var users []models.User
		if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if len(users) != 1 {
			t.Errorf("Expected a slice with one user, got %d users", len(users))
		}
	})

	t.Run("returns status 200 & empty array if no users are in db", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, _ := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/users", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var users []models.User
		if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected empty users slice, got %d users", len(users))
		}
	})
}
