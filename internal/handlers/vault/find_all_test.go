package vault_test

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

func TestFindAllVaults(t *testing.T) {
	t.Run("returns status 200 & array with vaults", func(t *testing.T) {
		serv, cleanup, db := testutils.NewTestApplication(t)
		defer cleanup()

		err := repositories.NewVaultRepo(db).CreateOne(context.Background(), "vault")
		if err != nil {
			t.Fatalf("failed to create new vault in repo: %v", err)
		}

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/vaults", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var vaults []models.Vault
		if err := json.NewDecoder(response.Body).Decode(&vaults); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if len(vaults) != 1 {
			t.Errorf("Expected a slice with one vault, got %d vaults", len(vaults))
		}
	})

	t.Run("returns status 200 & empty array if no vaults are in db", func(t *testing.T) {
		serv, cleanup, _ := testutils.NewTestApplication(t)
		defer cleanup()

		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/vaults", nil)
		serv.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		var vaults []models.Vault
		if err := json.NewDecoder(response.Body).Decode(&vaults); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if len(vaults) != 0 {
			t.Errorf("Expected empty vaults slice, got %d vaults", len(vaults))
		}
	})
}
