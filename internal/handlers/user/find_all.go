package user

import (
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/services"
	"github.com/kkstas/tnr-backend/internal/utils"
)

func FindAllHandler(logger *slog.Logger, userService *services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundUsers, err := userService.FindAll(r.Context())
		if err != nil {
			logger.Error("failed to find all users", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, foundUsers)
		if err != nil {
			logger.Error("failed to encode found users", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
