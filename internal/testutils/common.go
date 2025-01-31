package testutils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kkstas/tnr-backend/internal/app"
	_ "modernc.org/sqlite"
)

func NewTestApplication(t testing.TB) (newApp http.Handler, cleanup func(), db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	db, cleanupDb := openTestDB(t)

	cleanup = func() {
		cleanupDb()
		cancel()
	}

	newApp, err := app.NewApplication(ctx, db, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})))
	if err != nil {
		t.Fatalf("failed to create new application: %v", err)
	}

	return newApp, cleanup, db
}

func openTestDB(t testing.TB) (db *sql.DB, cleanup func()) {
	dbName := fmt.Sprintf("%s.db", RandomString(32))
	db, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(1)&_time_format=sqlite")
	if err != nil {
		t.Fatalf("failed to open sql db: %v", err)
	}

	cleanup = func() {
		db.Close()
		if err := os.Remove(dbName); err != nil {
			t.Fatalf("failed to remove test database file %s: %v", dbName, err)
		}
	}

	return db, cleanup
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func ToJSONBuffer(t testing.TB, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func DecodeJSON[T any](t testing.TB, body io.Reader) T {
	t.Helper()
	var result T
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	return result
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status: got %d, want %d", got, want)
	}
}

func AssertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
