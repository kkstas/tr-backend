package vault

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/services"
)

func FindAllHandler(logger *slog.Logger, vaultService *services.VaultService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundVaults, err := vaultService.FindAll(r.Context())
		if err != nil {
			logger.Error("err while finding all vaults", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(foundVaults)
		if err != nil {
			logger.Error("err while encoding found vaults", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
