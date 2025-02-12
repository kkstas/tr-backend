package session

import (
	"errors"
	"log/slog"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/utils"
)

var (
	minPasswordLength = 8
	maxPasswordLength = 50
	minNameLength     = 2
	maxNameLength     = 50
)

func RegisterHandler(logger *slog.Logger, userService *services.UserService) http.Handler {
	type reqBody struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[reqBody](r)
		if err != nil {
			utils.Encode(w, r, http.StatusBadRequest, map[string]string{"message": "failed to decode request body"})
			return

		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.FirstName, validation.Required, validation.Length(minNameLength, maxNameLength)),
			validation.Field(&body.LastName, validation.Required, validation.Length(minNameLength, maxNameLength)),
			validation.Field(&body.Email, validation.Required, is.EmailFormat),
			validation.Field(&body.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		)
		if err != nil {
			utils.Encode(w, r, http.StatusBadRequest, err)
			return
		}

		err = userService.CreateOne(r.Context(), body.FirstName, body.LastName, body.Email, body.Password)
		if err != nil {
			if errors.Is(err, services.ErrUserEmailAlreadyExists) {
				utils.Encode(w, r, http.StatusBadRequest, map[string]string{"email": "user with that email already exists"})
				return
			}
			logger.Error("failed to create user", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
