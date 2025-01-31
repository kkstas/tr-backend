package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/services"
)

func FindAllHandler(logger *slog.Logger, userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundUsers, err := userService.FindAll()
		if err != nil {
			logger.Error("err while finding all users", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(foundUsers)
		if err != nil {
			logger.Error("err while encoding found users", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
