package app

import (
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
	config *config.Config,
	db *sql.DB,
	logger *slog.Logger,
) *Application {

	app := new(Application)

	userRepo := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	vaultRepo := repositories.NewVaultRepo(db)
	vaultService := services.NewVaultService(vaultRepo, userService)
	expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
	expenseCategoryService := services.NewExpenseCategoryService(expenseCategoryRepo, vaultService)

	mux := handlers.SetupRoutes(config, logger, userService, vaultService, expenseCategoryService)
	app.Handler = middleware.LogHTTP(logger, mux)

	return app
}
