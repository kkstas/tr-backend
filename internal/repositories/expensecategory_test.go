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

func TestExpenseCategoryRepo_FindAll(t *testing.T) {
	t.Parallel()

	t.Run("finds all expense categories", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
		user := testutils.CreateTestUser(t, db)

		vaultID, err := repositories.NewVaultRepo(db).CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		_, err = expenseCategoryRepo.CreateOne(ctx, "category name", models.ExpenseCategoryStatusActive, 0, vaultID, user.ID)
		testutils.AssertNoError(t, err)

		foundCategories, err := expenseCategoryRepo.FindAll(ctx, vaultID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundCategories), 1)
	})

	t.Run("returns empty array if no categories are found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)

		foundCategories, err := expenseCategoryRepo.FindAll(ctx, uuid.New().String())
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundCategories), 0)
	})
}

func TestExpenseCategoryRepo_CreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates new expense category", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
		user := testutils.CreateTestUser(t, db)

		vaultID, err := repositories.NewVaultRepo(db).CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		name := "category name"
		status := models.ExpenseCategoryStatusActive
		priority := 0
		createdBy := user.ID

		categoryID, err := expenseCategoryRepo.CreateOne(ctx, name, status, priority, vaultID, createdBy)
		testutils.AssertNoError(t, err)

		category, err := expenseCategoryRepo.FindOneByID(ctx, categoryID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, category.Name, name)
		testutils.AssertEqual(t, category.Status, status)
		testutils.AssertEqual(t, category.Priority, priority)
		testutils.AssertEqual(t, category.VaultID, vaultID)
		testutils.AssertEqual(t, category.CreatedBy, createdBy)
	})
}

func TestExpenseCategoryRepo_FindOneByID(t *testing.T) {
	t.Parallel()

	t.Run("finds expense category by ID", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
		user := testutils.CreateTestUser(t, db)

		vaultID, err := repositories.NewVaultRepo(db).CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		name := "category name"
		status := models.ExpenseCategoryStatusActive
		priority := 0
		createdBy := user.ID

		categoryID, err := expenseCategoryRepo.CreateOne(ctx, name, status, priority, vaultID, createdBy)
		testutils.AssertNoError(t, err)

		foundCategory, err := expenseCategoryRepo.FindOneByID(ctx, categoryID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, foundCategory.Name, name)
		testutils.AssertEqual(t, foundCategory.Status, status)
		testutils.AssertEqual(t, foundCategory.Priority, priority)
		testutils.AssertEqual(t, foundCategory.VaultID, vaultID)
		testutils.AssertEqual(t, foundCategory.CreatedBy, createdBy)
	})

	t.Run("returns error if category is not found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)

		_, err := expenseCategoryRepo.FindOneByID(ctx, uuid.New().String())
		if err == nil {
			t.Error("expected an error but didn't get one")
		}
		want := repositories.ErrExpenseCategoryNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestExpenseCategoryRepo_SetStatus(t *testing.T) {
	t.Parallel()

	t.Run("sets expense category status", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
		user := testutils.CreateTestUser(t, db)

		vaultID, err := repositories.NewVaultRepo(db).CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		name := "category name"
		prevStatus := models.ExpenseCategoryStatusActive

		categoryID, err := expenseCategoryRepo.CreateOne(ctx, name, prevStatus, 0, vaultID, user.ID)
		testutils.AssertNoError(t, err)

		newStatus := models.ExpenseCategoryStatusInactive
		err = expenseCategoryRepo.SetStatus(ctx, categoryID, newStatus)
		testutils.AssertNoError(t, err)

		category, err := expenseCategoryRepo.FindOneByID(ctx, categoryID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, category.Status, newStatus)
	})
}

func TestExpenseCategoryRepo_SetPriority(t *testing.T) {
	t.Parallel()

	t.Run("sets expense category priority", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)
		user := testutils.CreateTestUser(t, db)

		vaultID, err := repositories.NewVaultRepo(db).CreateOne(ctx, user.ID, models.VaultRoleOwner, "vault name")
		testutils.AssertNoError(t, err)

		name := "category name"

		categoryID, err := expenseCategoryRepo.CreateOne(ctx, name, models.ExpenseCategoryStatusActive, 0, vaultID, user.ID)
		testutils.AssertNoError(t, err)

		newPriority := 1
		err = expenseCategoryRepo.SetPriority(ctx, categoryID, newPriority)
		testutils.AssertNoError(t, err)

		category, err := expenseCategoryRepo.FindOneByID(ctx, categoryID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, category.Priority, newPriority)
	})
}
