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
	"reflect"
	"testing"
	"time"

	_ "modernc.org/sqlite" // nolint: revive

	"github.com/kkstas/tr-backend/internal/app"
	"github.com/kkstas/tr-backend/internal/auth"
	"github.com/kkstas/tr-backend/internal/config"
	"github.com/kkstas/tr-backend/internal/database"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/services"
)

var jwtKey = []byte("secret-key")

func NewTestAppWithConfig(t testing.TB, config *config.Config) (newApp http.Handler, db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	db = OpenTestDB(t, ctx)
	t.Cleanup(cancel)

	// nolint: exhaustruct
	newApp = app.NewApplication(config, db, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})))

	return newApp, db

}

func NewTestApplication(t testing.TB) (newApp http.Handler, db *sql.DB) {
	config := &config.Config{
		EnableRegister: true,
		JWTSecretKey:   jwtKey,
	}
	return NewTestAppWithConfig(t, config)
}

func OpenTestDB(t testing.TB, ctx context.Context) (db *sql.DB) { // nolint: revive
	dbName := fmt.Sprintf("test-%s.db", RandomString(32))
	db, err := database.OpenDB(ctx, dbName)
	if err != nil {
		t.Fatalf("failed to open sql db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		if err := os.Remove(dbName); err != nil {
			t.Fatalf("failed to remove test database file %s: %v", dbName, err)
		}
	})

	return db
}

func NewTestUserService(db *sql.DB) *services.UserService {
	return services.NewUserService(repositories.NewUserRepo(db))
}

func NewTestVaultService(db *sql.DB) *services.VaultService {
	return services.NewVaultService(repositories.NewVaultRepo(db), NewTestUserService(db))
}

func NewTestExpenseCategoryService(db *sql.DB) *services.ExpenseCategoryService {
	return services.NewExpenseCategoryService(repositories.NewExpenseCategoryRepo(db), NewTestVaultService(db))
}

func CreateTestUser(t testing.TB, db *sql.DB) *models.User {
	userRepo := repositories.NewUserRepo(db)
	userEmail := RandomString(16) + "@email.com"
	err := userRepo.CreateOne(context.Background(), "firstName_"+RandomString(8), "lastName_"+RandomString(8), userEmail, "password")
	AssertNoError(t, err)

	createdUser, err := userRepo.FindOneByEmail(context.Background(), userEmail)
	AssertNoError(t, err)
	return createdUser
}

func CreateTestUserWithToken(t testing.TB, db *sql.DB) (token string, user *models.User) {
	createdUser := CreateTestUser(t, db)

	tkn, err := auth.CreateToken(jwtKey, createdUser.ID)
	AssertNoError(t, err)

	return tkn.Token, createdUser
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))] // nolint: gosec // using this instead of crypto/rand since it's only for testing
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

func AssertValidDate(t testing.TB, dateStr string) {
	t.Helper()
	layout := "2006-01-02T15:04:05Z"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		t.Errorf("string %s is not valid date in format %s: %v", dateStr, layout, err)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't expect one: %v", err)
	}
}

func AssertNotEmpty[T any](t testing.TB, got T) {
	t.Helper()

	v := reflect.ValueOf(got)

	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		t.Errorf("expected a non-nil value but got nil %T", got)
		return
	}

	switch v.Kind() { // nolint: exhaustive
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			t.Errorf("expected a non-empty %T but got empty", got)
		}
	case reflect.Ptr:
		return
	default:
		return
	}
}
