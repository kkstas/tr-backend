package vault

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

var (
	minVaultNameLength = 2
	maxVaultNameLength = 50
)

func CreateOne(
	vaultService *services.VaultService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	type reqBody struct {
		VaultName string `json:"vaultName"`
	}

	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		body, err := utils.Decode[reqBody](r)
		if err != nil {
			utils.Encode(w, r, http.StatusBadRequest, map[string]string{"message": "failed to decode request body"})
			return
		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.VaultName, validation.Required, validation.Length(minVaultNameLength, maxVaultNameLength)),
		)
		if err != nil {
			utils.Encode(w, r, http.StatusBadRequest, err)
			return
		}

		err = vaultService.CreateOne(r.Context(), user.ID, body.VaultName)
		if err != nil {
			utils.Encode(w, r, http.StatusInternalServerError, err.Error())
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
