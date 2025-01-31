package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/handlers/misc"
	"github.com/kkstas/tnr-backend/internal/handlers/user"
	"github.com/kkstas/tnr-backend/internal/services"
)

func SetupRoutes(logger *slog.Logger, db *sql.DB, userService *services.UserService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health-check", misc.HealthCheckHandler)
	mux.HandleFunc("/", misc.NotFoundHandler)

	mux.HandleFunc("GET /users", user.FindAllHandler(logger, userService))
	mux.HandleFunc("POST /users", user.CreateOneHandler(logger, userService))

	return mux
}
