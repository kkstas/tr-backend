package expensecategory

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

var (
	minCategoryNameLength = 2
	maxCategoryNameLength = 50
)

func CreateOne(
	expenseCategoryService *services.ExpenseCategoryService,
) func(w http.ResponseWriter, r *http.Request, user *models.User) {
	type reqBody struct {
		Name    string `json:"name"`
		VaultID string `json:"vaultID"`
	}

	return func(w http.ResponseWriter, r *http.Request, user *models.User) {
		body, err := utils.Decode[reqBody](r)
		if err != nil {
			utils.Encode(w, http.StatusBadRequest, map[string]string{"message": "failed to decode request body"})
			return
		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.Name, validation.Required, validation.Length(minCategoryNameLength, maxCategoryNameLength)),
			validation.Field(&body.VaultID, validation.Required),
		)
		if err != nil {
			utils.Encode(w, http.StatusBadRequest, err)
			return
		}

		err = expenseCategoryService.CreateOne(r.Context(), body.Name, user.ID, body.VaultID)
		if err != nil {
			if errors.Is(err, services.ErrVaultNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			utils.Encode(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
