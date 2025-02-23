package repositories_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/repositories"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestUserRepo_CreateOne(t *testing.T) {
	t.Parallel()

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
	})
}

func TestUserRepo_FindOneByID(t *testing.T) {
	t.Parallel()

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

		want := repositories.ErrUserNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestUserRepo_FindOneByEmail(t *testing.T) {
	t.Parallel()

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
		want := repositories.ErrUserNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestUserRepo_FindPasswordHashAndUserIDForEmail(t *testing.T) {
	t.Parallel()

	t.Run("finds password hash and user ID for email", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@email.com"
		passwordHash := "xyz"

		err := userRepo.CreateOne(ctx, firstName, lastName, email, passwordHash)
		testutils.AssertNoError(t, err)

		foundPasswordHash, _, err := userRepo.FindPasswordHashAndUserIDForEmail(ctx, email)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, foundPasswordHash, passwordHash)
	})

	t.Run("returns error if user with given email is not found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepo := repositories.NewUserRepo(testutils.OpenTestDB(t, ctx))

		_, _, err := userRepo.FindPasswordHashAndUserIDForEmail(ctx, "idontexist@email.com")
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := repositories.ErrUserNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestUserRepo_FindAll(t *testing.T) {
	t.Parallel()

	t.Run("finds all users", func(t *testing.T) {
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

		foundUser := foundUsers[0]

		testutils.AssertEqual(t, foundUser.FirstName, firstName)
		testutils.AssertEqual(t, foundUser.LastName, lastName)
		testutils.AssertEqual(t, foundUser.Email, email)
		testutils.AssertValidDate(t, foundUser.CreatedAt)
		if err := uuid.Validate(foundUser.ID); err != nil {
			t.Errorf("expected id to be valid uuid, got error: %v", err)
		}
	})

	t.Run("returns empty array if no users are found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		foundUsers, err := userRepo.FindAll(ctx)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundUsers), 0)
	})
}

func TestUserRepo_AssignActiveVault(t *testing.T) {
	t.Parallel()

	t.Run("assigns active vault to user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)
		vaultRepo := repositories.NewVaultRepo(db)

		email := "john.doe@email.com"

		err := userRepo.CreateOne(ctx, "John", "Doe", email, "somepassword")
		testutils.AssertNoError(t, err)

		user, err := userRepo.FindOneByEmail(ctx, email)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, user.ActiveVault, "")

		vaultID, err := vaultRepo.CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		err = userRepo.AssignActiveVault(ctx, user.ID, vaultID)
		testutils.AssertNoError(t, err)

		user, err = userRepo.FindOneByEmail(ctx, email)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, user.ActiveVault, vaultID)
	})

	t.Run("returns error if vault does not exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userRepo := repositories.NewUserRepo(db)

		email := "john.doe@email.com"

		err := userRepo.CreateOne(ctx, "John", "Doe", email, "somepassword")
		testutils.AssertNoError(t, err)

		user, err := userRepo.FindOneByEmail(ctx, email)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, user.ActiveVault, "")

		err = userRepo.AssignActiveVault(ctx, user.ID, uuid.New().String())
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		user, err = userRepo.FindOneByEmail(ctx, email)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, user.ActiveVault, "")
	})
}
