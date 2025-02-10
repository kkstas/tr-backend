package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/config"
	"github.com/kkstas/tnr-backend/internal/handlers"
	"github.com/kkstas/tnr-backend/internal/middleware"
	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/services"
)

type Application struct {
	http.Handler
}

func NewApplication(
	ctx context.Context,
	config *config.Config,
	db *sql.DB,
	logger *slog.Logger,
) *Application {

	app := new(Application)

	userRepo := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	vaultRepo := repositories.NewVaultRepo(db)
	vaultService := services.NewVaultService(vaultRepo)

	mux := handlers.SetupRoutes(config, logger, db, userService, vaultService)

	app.Handler = middleware.LogHttp(logger, mux)

	return app
}
