package vault

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

func AddUser(
	vaultService *services.VaultService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	type reqBody struct {
		UserID string `json:"userID"`
		Role   string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		vaultID := r.PathValue("vaultID")

		body, err := utils.Decode[reqBody](r)
		if err != nil {
			utils.Encode(w, http.StatusBadRequest, map[string]string{"message": "failed to decode request body"})
			return
		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.UserID, validation.Required),
			validation.Field(&body.Role, validation.Required, validation.In(string(models.VaultRoleEditor))),
		)
		if err != nil {
			utils.Encode(w, http.StatusBadRequest, err)
			return
		}

		err = vaultService.AddUser(r.Context(), user.ID, body.UserID, vaultID, models.VaultRole(body.Role))
		if err != nil {
			if errors.Is(err, services.ErrUserAlreadyAssignedToVault) {
				utils.Encode(w, http.StatusBadRequest, map[string]string{"message": "user has already been assigned to this vault"})
				return
			}

			if errors.Is(err, services.ErrVaultNotFound) {
				utils.Encode(w, http.StatusNotFound, map[string]string{"message": "vault not found"})
				return
			}
			utils.Encode(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
