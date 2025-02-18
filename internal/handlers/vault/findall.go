package vault

import (
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

func FindAll(vaultService *services.VaultService) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		vaults, err := vaultService.FindAll(r.Context(), user.ID)
		if err != nil {
			utils.Encode(w, http.StatusInternalServerError, err.Error())
		}

		utils.Encode(w, http.StatusOK, vaults)
	}
}
