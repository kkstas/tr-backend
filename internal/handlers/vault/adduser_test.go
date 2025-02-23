package vault_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestAddUserToVault(t *testing.T) {
	t.Parallel()

	t.Run("adds user to vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		inviterToken, _ := testutils.CreateTestUserWithToken(t, db)
		inviteeToken, invitee := testutils.CreateTestUserWithToken(t, db)

		var createdVaultID string

		{
			// create vault
			vaultFC := struct {
				VaultName string `json:"vaultName"`
			}{VaultName: "asdf"}

			request := httptest.NewRequest("POST", "/vaults", testutils.ToJSONBuffer(t, vaultFC))
			request.Header.Set("Authorization", "Bearer "+inviterToken)
			response := httptest.NewRecorder()
			serv.ServeHTTP(response, request)

			testutils.AssertStatus(t, response.Code, http.StatusNoContent)

			// get vault ID
			request = httptest.NewRequest("GET", "/vaults", nil)
			request.Header.Set("Authorization", "Bearer "+inviterToken)
			response = httptest.NewRecorder()
			serv.ServeHTTP(response, request)
			var vaults = []models.UserVaultWithRole{}
			err := json.NewDecoder(response.Body).Decode(&vaults)
			testutils.AssertNoError(t, err)
			createdVaultID = vaults[0].ID
		}

		invitationReqBody := struct {
			UserID string           `json:"userID"`
			Role   models.VaultRole `json:"role"`
		}{UserID: invitee.ID, Role: models.VaultRoleEditor}

		request := httptest.NewRequest(
			"POST",
			fmt.Sprintf("/vaults/%s/users", createdVaultID),
			testutils.ToJSONBuffer(t, invitationReqBody),
		)
		request.Header.Set("Authorization", "Bearer "+inviterToken)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		// check if invited user has been added to vault
		request = httptest.NewRequest("GET", "/vaults", nil)
		request.Header.Set("Authorization", "Bearer "+inviteeToken)
		response = httptest.NewRecorder()
		serv.ServeHTTP(response, request)
		var vaults = []models.UserVaultWithRole{}
		err := json.NewDecoder(response.Body).Decode(&vaults)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, vaults[0].ID, createdVaultID)
	})

	t.Run("returns 404 if vault does not exist", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		inviterToken, _ := testutils.CreateTestUserWithToken(t, db)
		_, invitee := testutils.CreateTestUserWithToken(t, db)

		invitationReqBody := struct {
			UserID string           `json:"userID"`
			Role   models.VaultRole `json:"role"`
		}{UserID: invitee.ID, Role: models.VaultRoleEditor}

		request := httptest.NewRequest(
			"POST",
			fmt.Sprintf("/vaults/%s/users", uuid.New().String()),
			testutils.ToJSONBuffer(t, invitationReqBody),
		)
		request.Header.Set("Authorization", "Bearer "+inviterToken)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 400 if user has already been assigned to this vault", func(t *testing.T) {
		t.Parallel()
		serv, db := testutils.NewTestApplication(t)
		inviterToken, _ := testutils.CreateTestUserWithToken(t, db)
		_, invitee := testutils.CreateTestUserWithToken(t, db)

		var createdVaultID string

		{
			// create vault
			vaultFC := struct {
				VaultName string `json:"vaultName"`
			}{VaultName: "asdf"}

			request := httptest.NewRequest("POST", "/vaults", testutils.ToJSONBuffer(t, vaultFC))
			request.Header.Set("Authorization", "Bearer "+inviterToken)
			response := httptest.NewRecorder()
			serv.ServeHTTP(response, request)

			testutils.AssertStatus(t, response.Code, http.StatusNoContent)

			// get vault ID
			request = httptest.NewRequest("GET", "/vaults", nil)
			request.Header.Set("Authorization", "Bearer "+inviterToken)
			response = httptest.NewRecorder()
			serv.ServeHTTP(response, request)
			var vaults = []models.UserVaultWithRole{}
			err := json.NewDecoder(response.Body).Decode(&vaults)
			testutils.AssertNoError(t, err)
			createdVaultID = vaults[0].ID
		}

		invitationReqBody := struct {
			UserID string           `json:"userID"`
			Role   models.VaultRole `json:"role"`
		}{UserID: invitee.ID, Role: models.VaultRoleEditor}

		{ // add user to vault
			request := httptest.NewRequest(
				"POST",
				fmt.Sprintf("/vaults/%s/users", createdVaultID),
				testutils.ToJSONBuffer(t, invitationReqBody),
			)
			request.Header.Set("Authorization", "Bearer "+inviterToken)
			response := httptest.NewRecorder()
			serv.ServeHTTP(response, request)

			testutils.AssertStatus(t, response.Code, http.StatusNoContent)
		}

		request := httptest.NewRequest(
			"POST",
			fmt.Sprintf("/vaults/%s/users", createdVaultID),
			testutils.ToJSONBuffer(t, invitationReqBody),
		)
		request.Header.Set("Authorization", "Bearer "+inviterToken)
		response := httptest.NewRecorder()
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

	})
}
