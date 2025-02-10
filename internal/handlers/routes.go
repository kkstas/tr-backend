package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/handlers/misc"
	"github.com/kkstas/tnr-backend/internal/handlers/session"
	"github.com/kkstas/tnr-backend/internal/handlers/user"
	"github.com/kkstas/tnr-backend/internal/handlers/vault"
	mw "github.com/kkstas/tnr-backend/internal/middleware"
	"github.com/kkstas/tnr-backend/internal/services"
)

func SetupRoutes(
	logger *slog.Logger,
	db *sql.DB,
	userService *services.UserService,
	vaultService *services.VaultService,
	enableRegister bool,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health-check", misc.HealthCheckHandler)
	mux.HandleFunc("/", misc.NotFoundHandler)

	mux.Handle("POST /login", session.LoginHandler(logger, userService))
	mux.Handle("POST /register", mw.Enable(enableRegister, session.RegisterHandler(logger, userService)))

	mux.Handle("GET /users", user.FindAllHandler(logger, userService))
	mux.Handle("GET /users/{id}", user.FindOneByIDHandler(logger, userService))
	mux.Handle("POST /vaults", vault.CreateOneHandler(logger, vaultService))
	mux.Handle("GET /vaults", vault.FindAllHandler(logger, vaultService))

	return mux
}
