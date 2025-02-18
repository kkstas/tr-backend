package middleware

import (
	"log/slog"
	"net/http"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
)

func WithUser(
	logger *slog.Logger,
	userService *services.UserService,
) func(fn func(w http.ResponseWriter, r *http.Request, user *models.User)) http.HandlerFunc {
	return func(fn func(w http.ResponseWriter, r *http.Request, user *models.User)) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(JWTClaims)
			if !ok {
				logger.Error("user claims not found in request context")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := userService.FindOneByID(r.Context(), claims.UserID)
			if err != nil {
				logger.Error("failed to find user with id from token %s: %v", claims.UserID, err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fn(w, r, user)
		})
	}
}
