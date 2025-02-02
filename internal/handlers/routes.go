package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/handlers/misc"
	"github.com/kkstas/tnr-backend/internal/handlers/session"
	"github.com/kkstas/tnr-backend/internal/handlers/user"
	"github.com/kkstas/tnr-backend/internal/handlers/vault"
	"github.com/kkstas/tnr-backend/internal/services"
)

func SetupRoutes(logger *slog.Logger, db *sql.DB, userService *services.UserService, vaultService *services.VaultService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health-check", misc.HealthCheckHandler)
	mux.HandleFunc("/", misc.NotFoundHandler)

	mux.Handle("POST /login", session.LoginHandler(logger, userService))
	mux.Handle("POST /register", session.RegisterHandler(logger, userService))

	mux.Handle("GET /users", user.FindAllHandler(logger, userService))
	mux.Handle("GET /users/{id}", user.FindOneByIDHandler(logger, userService))
	mux.Handle("POST /vaults", vault.CreateOneHandler(logger, vaultService))
	mux.Handle("GET /vaults", vault.FindAllHandler(logger, vaultService))

	return mux
}
