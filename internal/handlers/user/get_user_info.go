package user

import (
	"log/slog"
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

func GetUserInfo(
	logger *slog.Logger,
	userService *services.UserService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		err := utils.Encode(w, r, http.StatusOK, user)
		if err != nil {
			logger.Error("failed to encode vaults", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
