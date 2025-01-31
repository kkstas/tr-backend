package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/models"
	"github.com/kkstas/tnr-backend/internal/services"
)

func CreateOneHandler(logger *slog.Logger, userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Error("failed to decode request body", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = userService.CreateOne(user.FirstName, user.LastName, user.Email)
		if err != nil {
			logger.Error("failed to create new user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
