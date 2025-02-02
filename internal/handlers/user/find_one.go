package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kkstas/tnr-backend/internal/services"
	"github.com/kkstas/tnr-backend/internal/utils"
)

func FindOneByIDHandler(logger *slog.Logger, userService *services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := uuid.Validate(id); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"id":"invalid uuid"}`)
			return
		}

		foundUser, err := userService.FindOneByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, foundUser)
		if err != nil {
			logger.Error("failed to encode found user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
