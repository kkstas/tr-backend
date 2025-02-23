package repositories_test

import (
	"context"
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

		err = expenseCategoryRepo.CreateOne(ctx, "category name", models.ExpenseCategoryStatusActive, 0, vaultID, user.ID)
		testutils.AssertNoError(t, err)

		foundVaults, err := expenseCategoryRepo.FindAll(ctx, vaultID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)
	})

	t.Run("returns empty array if no categories are found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		expenseCategoryRepo := repositories.NewExpenseCategoryRepo(db)

		foundVaults, err := expenseCategoryRepo.FindAll(ctx, uuid.New().String())
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 0)
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

		err = expenseCategoryRepo.CreateOne(ctx, name, status, priority, vaultID, createdBy)
		testutils.AssertNoError(t, err)

		foundVaults, err := expenseCategoryRepo.FindAll(ctx, vaultID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)

		vault := foundVaults[0]

		testutils.AssertEqual(t, vault.Name, name)
		testutils.AssertEqual(t, vault.Status, status)
		testutils.AssertEqual(t, vault.Priority, priority)
		testutils.AssertEqual(t, vault.VaultID, vaultID)
		testutils.AssertEqual(t, vault.CreatedBy, createdBy)
	})
}
