package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestCreateOne(t *testing.T) {
	t.Run("creates one user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userService := services.NewUserService(repositories.NewUserRepo(db))

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@email.com"

		err := userService.CreateOne(ctx, firstName, lastName, email, "somepassword")
		testutils.AssertNoError(t, err)

		foundUsers, err := userService.FindAll(ctx)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(foundUsers), 1)
		testutils.AssertEqual(t, foundUsers[0].FirstName, firstName)
		testutils.AssertEqual(t, foundUsers[0].LastName, lastName)
		testutils.AssertEqual(t, foundUsers[0].Email, email)
		testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
		if err := uuid.Validate(foundUsers[0].ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})

	t.Run("returns error if user with given email already exists", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userService := services.NewUserService(repositories.NewUserRepo(db))

		email := "some@email.com"

		err := userService.CreateOne(ctx, "John", "Doe", email, "somepassword")
		testutils.AssertNoError(t, err)

		err = userService.CreateOne(ctx, "John", "Doe", email, "somepassword")
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		if !errors.Is(err, services.ErrUserEmailAlreadyExists) {
			t.Errorf("expected error '%v', got '%v'", services.ErrUserEmailAlreadyExists, err)
		}
	})
}
