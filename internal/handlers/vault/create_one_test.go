package vault_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestCreateOneVault(t *testing.T) {
	t.Run("returns status 204 & saves vault in DB when created new vault", func(t *testing.T) {
		t.Parallel()
		serv, _ := testutils.NewTestApplication(t)

		vaultFC := models.Vault{Name: "Doe"}

		reqBody := testutils.ToJSONBuffer(t, vaultFC)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/vaults", reqBody)
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		// Check if data has been saved properly
		response = httptest.NewRecorder()
		serv.ServeHTTP(response, httptest.NewRequest("GET", "/vaults", nil))
		foundVaults := testutils.DecodeJSON[[]models.Vault](t, response.Body)
		want := 1
		if len(foundVaults) != want {
			t.Errorf("got %d vaults, want %d", len(foundVaults), want)
		}

		testutils.AssertEqual(t, foundVaults[0].Name, vaultFC.Name)
		testutils.AssertValidDate(t, foundVaults[0].CreatedAt)
		if err := uuid.Validate(foundVaults[0].ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})
}
