package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kkstas/tnr-backend/internal/database"
	"github.com/kkstas/tnr-backend/internal/handlers"
	"github.com/kkstas/tnr-backend/internal/middleware"
	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/services"
)

type Application struct {
	http.Handler
}

func NewApplication(ctx context.Context, db *sql.DB, logger *slog.Logger) (*Application, error) {
	err := database.InitDBTables(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to init database tables: %w", err)
	}

	app := new(Application)

	userRepo := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepo)

	mux := handlers.SetupRoutes(logger, db, userService)

	app.Handler = middleware.LogHttp(logger, mux)

	return app, nil
}
