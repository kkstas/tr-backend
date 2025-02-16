package vault

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
)

func DeleteOneByID(
	logger *slog.Logger,
	vaultService *services.VaultService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {

	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		vaultID := r.PathValue("id")
		if vaultID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := vaultService.DeleteOneByID(r.Context(), user.ID, vaultID)
		if err != nil {
			if errors.Is(err, services.ErrVaultNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if errors.Is(err, services.ErrInsufficientVaultPermissions) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			logger.Error("failed to delete vault", "vaultID", vaultID, "userID", user.ID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
