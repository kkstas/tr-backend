package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/kkstas/tr-backend/internal/config"
	"github.com/kkstas/tr-backend/internal/handlers"
	"github.com/kkstas/tr-backend/internal/middleware"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/services"
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
	vaultService := services.NewVaultService(vaultRepo, userService)

	mux := handlers.SetupRoutes(config, logger, db, userService, vaultService)

	app.Handler = middleware.LogHttp(logger, mux)

	return app
}
