package vault

import (
	"fmt"
	"log/slog"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

func CreateOneHandler(logger *slog.Logger, vaultService *services.VaultService) http.Handler {
	type reqBody struct {
		Name string `json:"name"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[reqBody](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.Name, validation.Required, validation.Length(2, 50)),
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, utils.ToJSON(err))
			return
		}

		err = vaultService.CreateOne(r.Context(), body.Name)
		if err != nil {
			logger.Error("failed to create new vault", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
