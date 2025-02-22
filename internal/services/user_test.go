package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestUserService_CreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates one user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userService := testutils.NewTestUserService(db)

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
		userService := testutils.NewTestUserService(db)

		email := "some@email.com"

		err := userService.CreateOne(ctx, "John", "Doe", email, "somepassword")
		testutils.AssertNoError(t, err)

		err = userService.CreateOne(ctx, "John", "Doe", email, "somepassword")
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		want := services.ErrUserEmailAlreadyExists
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}
