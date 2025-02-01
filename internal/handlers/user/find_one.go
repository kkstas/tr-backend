package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kkstas/tnr-backend/internal/services"
)

func FindOneByIDHandler(logger *slog.Logger, userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = json.NewEncoder(w).Encode(foundUser)
		if err != nil {
			logger.Error("err while encoding found user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
