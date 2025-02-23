package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestExpenseCategoryService_CreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates expense category", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)
		_, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		err := expenseCategoryService.CreateOne(ctx, "category name", user.ID, vault.ID)
		if err != nil {
			t.Errorf("didn't expect an error, but got one: %v", err)
		}
	})

	t.Run("returns error if user is not related to that vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)

		_, _, vault := testutils.CreateTestUserWithTokenAndVault(t, db)
		user := testutils.CreateTestUser(t, db)

		err := expenseCategoryService.CreateOne(ctx, "category name", user.ID, vault.ID)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})

	t.Run("returns error if user is related to that vault but is not an owner", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)
		vaultService := testutils.NewTestVaultService(db)

		_, vaultOwner, vault := testutils.CreateTestUserWithTokenAndVault(t, db)
		user := testutils.CreateTestUser(t, db)

		err := vaultService.AddUser(ctx, vaultOwner.ID, user.ID, vault.ID, models.VaultRoleEditor)
		testutils.AssertNoError(t, err)

		err = expenseCategoryService.CreateOne(ctx, "category name", user.ID, vault.ID)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}
		want := services.ErrUserIsNotVaultOwner
		if !errors.Is(err, services.ErrUserIsNotVaultOwner) {
			t.Errorf("expected %q, got %v", want, err)
		}
	})

	t.Run("returns error when expense category with given name already exists for this vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)
		_, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		err := expenseCategoryService.CreateOne(ctx, "category name", user.ID, vault.ID)
		testutils.AssertNoError(t, err)

		err = expenseCategoryService.CreateOne(ctx, "category name", user.ID, vault.ID)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := services.ErrExpenseCategoryWithThatNameAlreadyExists
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestExpenseCategoryService_FindAll(t *testing.T) {
	t.Parallel()

	t.Run("finds all expense categories", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)
		_, user, vault := testutils.CreateTestUserWithTokenAndVault(t, db)

		err := expenseCategoryService.CreateOne(ctx, "category one", user.ID, vault.ID)
		testutils.AssertNoError(t, err)
		err = expenseCategoryService.CreateOne(ctx, "category two", user.ID, vault.ID)
		testutils.AssertNoError(t, err)

		foundCategories, err := expenseCategoryService.FindAll(ctx, user.ID, vault.ID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, len(foundCategories), 2)

		for _, category := range foundCategories {
			testutils.AssertEqual(t, category.CreatedBy, user.ID)
			testutils.AssertEqual(t, category.VaultID, vault.ID)
			testutils.AssertEqual(t, category.Status, models.ExpenseCategoryStatusActive)
			testutils.AssertEqual(t, category.Priority, 0)
		}
	})

	t.Run("returns error if user does not belong to provided vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryService := testutils.NewTestExpenseCategoryService(db)

		_, vaultOwner, vault := testutils.CreateTestUserWithTokenAndVault(t, db)
		user := testutils.CreateTestUser(t, db)

		err := expenseCategoryService.CreateOne(ctx, "category name", vaultOwner.ID, vault.ID)
		testutils.AssertNoError(t, err)

		_, err = expenseCategoryService.FindAll(ctx, user.ID, vault.ID)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
