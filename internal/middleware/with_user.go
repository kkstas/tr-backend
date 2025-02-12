package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kkstas/tr-backend/internal/auth"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
)

func WithUser(jwtSecretKey []byte, logger *slog.Logger, userService *services.UserService) func(fn func(w http.ResponseWriter, r *http.Request, user *models.User)) http.HandlerFunc {
	return func(fn func(w http.ResponseWriter, r *http.Request, user *models.User)) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if !strings.HasPrefix(tokenString, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := auth.VerifyToken(jwtSecretKey, tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				if !errors.Is(err, auth.ErrInvalidToken) {
					logger.Error("failed to verify token", "error", err)
				}
				return
			}

			id, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := userService.FindOneByID(r.Context(), id)
			if err != nil {
				logger.Error("failed to find user with id from token %s: %v", id, err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fn(w, r, user)
		})
	}
}
