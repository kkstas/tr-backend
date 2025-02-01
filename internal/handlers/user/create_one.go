package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/kkstas/tnr-backend/internal/services"
	"github.com/kkstas/tnr-backend/internal/utils"
)

func CreateOneHandler(logger *slog.Logger, userService *services.UserService) http.HandlerFunc {
	type reqBody struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body reqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			logger.Error("failed to decode request body", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.FirstName, validation.Required, validation.Length(2, 50)),
			validation.Field(&body.LastName, validation.Required, validation.Length(2, 50)),
			validation.Field(&body.Email, validation.Required, is.EmailFormat),
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, utils.ToJSON(err))
			return
		}

		err = userService.CreateOne(r.Context(), body.FirstName, body.LastName, body.Email)
		if err != nil {
			logger.Error("failed to create new user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
