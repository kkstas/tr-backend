package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/config"
	"github.com/kkstas/tnr-backend/internal/handlers/misc"
	"github.com/kkstas/tnr-backend/internal/handlers/session"
	"github.com/kkstas/tnr-backend/internal/handlers/user"
	"github.com/kkstas/tnr-backend/internal/handlers/vault"
	mw "github.com/kkstas/tnr-backend/internal/middleware"
	"github.com/kkstas/tnr-backend/internal/services"
)

func SetupRoutes(
	cfg *config.Config,
	logger *slog.Logger,
	db *sql.DB,
	userService *services.UserService,
	vaultService *services.VaultService,
) http.Handler {
	mux := http.NewServeMux()

	withUser := mw.WithUser(cfg.JWTSecretKey, logger, userService)

	mux.HandleFunc("GET /health-check", misc.HealthCheckHandler)
	mux.HandleFunc("/", misc.NotFoundHandler)

	mux.Handle("POST /login", session.LoginHandler(cfg.JWTSecretKey, logger, userService))
	mux.Handle("POST /register", mw.Enable(cfg.EnableRegister, session.RegisterHandler(logger, userService)))

	mux.Handle("GET /user", withUser(user.GetUserInfo(logger, userService)))
	mux.Handle("GET /users/{id}", user.FindOneByIDHandler(logger, userService))
	mux.Handle("POST /vaults", vault.CreateOneHandler(logger, vaultService))
	mux.Handle("GET /vaults", vault.FindAllHandler(logger, vaultService))

	return mux
}
