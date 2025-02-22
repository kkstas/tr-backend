package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kkstas/tr-backend/internal/auth"
)

type UserClaimsKeyType struct{}

var UserClaimsKey = UserClaimsKeyType{}

type JWTClaims struct {
	UserID string
}

func RequireAuth(
	jwtSecretKey []byte,
	logger *slog.Logger,
) func(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
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

			jwtClaims := JWTClaims{} // nolint: exhaustruct

			claims := token.Claims.(jwt.MapClaims)

			if val, ok := claims["sub"].(string); ok {
				jwtClaims.UserID = val
			} else {
				logger.Error("failed to find UserID in claims", "claims", jwtClaims)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, jwtClaims)
			fn(w, r.WithContext(ctx))
		})
	}
}
