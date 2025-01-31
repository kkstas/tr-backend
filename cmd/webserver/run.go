package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kkstas/tnr-backend/internal/app"
	"github.com/kkstas/tnr-backend/internal/database"
	_ "modernc.org/sqlite"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	db, err := openDB()
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()
	err = database.InitDBTables(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to init db tables: %w", err)
	}

	app := app.NewApplication(ctx, db, initLogger(os.Stdout))

	server := &http.Server{
		Addr:              ":8000",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           app,
	}

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to ListenAndServe: %w", err)
	}

	return nil
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./database.db?_pragma=foreign_keys(1)&_time_format=sqlite")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

func initLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	))
}
