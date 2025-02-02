package vault

import (
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/services"
	"github.com/kkstas/tnr-backend/internal/utils"
)

func FindAllHandler(logger *slog.Logger, vaultService *services.VaultService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundVaults, err := vaultService.FindAll(r.Context())
		if err != nil {
			logger.Error("failed to find all vaults", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, foundVaults)
		if err != nil {
			logger.Error("failed to encode vaults", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
