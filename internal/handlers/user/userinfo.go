package user

import (
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/utils"
)

func GetUserInfo() func(w http.ResponseWriter, r *http.Request, user *models.User) {
	return func(w http.ResponseWriter, _ *http.Request, user *models.User) {
		utils.Encode(w, http.StatusOK, user)
	}
}
