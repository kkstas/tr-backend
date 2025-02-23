package expensecategory

import (
	"errors"
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

func FindAll(
	expenseCategoryService *services.ExpenseCategoryService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		vaultID := r.PathValue("vaultID")

		categories, err := expenseCategoryService.FindAll(r.Context(), user.ID, vaultID)
		if err != nil {
			if errors.Is(err, services.ErrVaultNotFound) {
				utils.Encode(w, http.StatusNotFound, map[string]string{"message": "vault not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.Encode(w, http.StatusOK, categories)
	}
}
