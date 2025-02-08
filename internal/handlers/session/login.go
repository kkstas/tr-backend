package session

import (
	"errors"
	"log/slog"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/kkstas/tnr-backend/internal/auth"
	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/services"
	"github.com/kkstas/tnr-backend/internal/utils"
)

func LoginHandler(logger *slog.Logger, userService *services.UserService) http.Handler {
	type loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := utils.Decode[loginData](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		err = validation.ValidateStruct(&body,
			validation.Field(&body.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
			validation.Field(&body.Email, validation.Required, is.EmailFormat),
		)
		if err != nil {
			_ = utils.Encode(w, r, http.StatusBadRequest, err)
			return
		}

		passwordHash, userID, err := userService.FindPasswordHashAndUserIDForEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, repositories.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error("failed to find password hash and user ID for email", "email", body.Email, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !auth.CheckPassword(passwordHash, body.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := auth.CreateToken(userID)
		if err != nil {
			logger.Error("failed to create token", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = utils.Encode(w, r, http.StatusOK, token)
	})
}
