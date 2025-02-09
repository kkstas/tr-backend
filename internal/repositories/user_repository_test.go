package repositories_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateOne(t *testing.T) {
	t.Run("creates one user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@email.com"

		err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
		testutils.AssertNoError(t, err)

		foundUsers, err := userRepo.FindAll(ctx)
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
		userRepo := repositories.NewUserRepo(db)

		email := "some@email.com"

		err := userRepo.CreateOne(ctx, "John", "Doe", email, "somepassword")
		testutils.AssertNoError(t, err)

		err = userRepo.CreateOne(ctx, "John", "Doe", email, "somepassword")
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		if !errors.Is(err, repositories.ErrUserEmailAlreadyExists) {
			t.Errorf("expected error '%v', got '%v'", repositories.ErrUserEmailAlreadyExists, err)
		}
	})
}

func TestFindOneUserByID(t *testing.T) {
	t.Run("finds one user by ID", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@email.com"

		err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
		testutils.AssertNoError(t, err)

		foundUsers, err := userRepo.FindAll(ctx)
		testutils.AssertNoError(t, err)

		foundUser, err := userRepo.FindOneByID(ctx, foundUsers[0].ID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, foundUser.FirstName, firstName)
		testutils.AssertEqual(t, foundUser.LastName, lastName)
		testutils.AssertEqual(t, foundUser.Email, email)
		testutils.AssertValidDate(t, foundUser.CreatedAt)
		if err := uuid.Validate(foundUser.ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})

	t.Run("returns correct error if user does not exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		_, err := userRepo.FindOneByID(ctx, uuid.New().String())
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		if !errors.Is(err, repositories.ErrUserNotFound) {
			t.Errorf("expected error '%v', got '%v'", repositories.ErrUserNotFound, err)
		}
	})
}

func TestFindOneUserByEmail(t *testing.T) {
	t.Run("finds one user by email", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@email.com"

		err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
		testutils.AssertNoError(t, err)

		foundUser, err := userRepo.FindOneByEmail(ctx, email)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, foundUser.FirstName, firstName)
		testutils.AssertEqual(t, foundUser.LastName, lastName)
		testutils.AssertEqual(t, foundUser.Email, email)
		testutils.AssertValidDate(t, foundUser.CreatedAt)
		if err := uuid.Validate(foundUser.ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})

	t.Run("returns correct error if user does not exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		_, err := userRepo.FindOneByEmail(ctx, "asdf@asdf.com")
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		if !errors.Is(err, repositories.ErrUserNotFound) {
			t.Errorf("expected error '%v', got '%v'", repositories.ErrUserNotFound, err)
		}
	})
}
