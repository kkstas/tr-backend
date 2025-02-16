package services_test

import (
	"context"
	"testing"

	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestFindOneByID(t *testing.T) {
	t.Parallel()

	t.Run("finds one vault for user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultService := testutils.NewTestVaultService(db)
		_, createdUser := testutils.CreateUserWithToken(t, db)

		vaultName := "some vault"

		err := vaultService.CreateOne(ctx, createdUser.ID, vaultName)
		testutils.AssertNoError(t, err)

		vaults, err := vaultService.FindAll(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)

		foundVault, err := vaultService.FindOneByID(ctx, createdUser.ID, vaults[0].ID)
		testutils.AssertNoError(t, err)

		testutils.AssertNotEmpty(t, foundVault)
		testutils.AssertEqual(t, foundVault.ID, vaults[0].ID)
		testutils.AssertEqual(t, foundVault.Name, vaultName)
		testutils.AssertEqual(t, foundVault.UserRole, models.VaultRoleOwner)
	})
}

func TestVaultService_CreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates new vault and assigns vault ID to user's active_vault if active_vault was empty", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		userService := testutils.NewTestUserService(db)
		vaultService := testutils.NewTestVaultService(db)
		_, createdUser := testutils.CreateUserWithToken(t, db)

		vaultName := "some vault"

		err := vaultService.CreateOne(ctx, createdUser.ID, vaultName)
		testutils.AssertNoError(t, err)

		vaults, err := vaultService.FindAll(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)

		foundVault := vaults[0]

		testutils.AssertNotEmpty(t, foundVault)
		testutils.AssertEqual(t, foundVault.ID, vaults[0].ID)
		testutils.AssertEqual(t, foundVault.Name, vaultName)
		testutils.AssertEqual(t, foundVault.UserRole, models.VaultRoleOwner)

		newUserData, err := userService.FindOneByID(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)

		if newUserData.ActiveVault != foundVault.ID {
			t.Errorf("expected new vault ID to be assigned as user's active vault, got %s", newUserData.ActiveVault)
		}
	})
}

func TestVaultService_DeleteOneByID(t *testing.T) {
	t.Parallel()

	t.Run("deletes vault and clears user's active_vault property if deleted vault was the active one", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		_, createdUser := testutils.CreateUserWithToken(t, db)
		userService := testutils.NewTestUserService(db)
		vaultService := testutils.NewTestVaultService(db)

		err := vaultService.CreateOne(ctx, createdUser.ID, "vault name")
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultService.FindAll(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)

		userData, err := userService.FindOneByID(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, userData.ActiveVault, foundVaults[0].ID)

		err = vaultService.DeleteOneByID(ctx, createdUser.ID, foundVaults[0].ID)
		testutils.AssertNoError(t, err)

		foundVaults, err = vaultService.FindAll(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 0)

		newUserData, err := userService.FindOneByID(ctx, createdUser.ID)
		testutils.AssertNoError(t, err)

		if newUserData.ActiveVault != "" {
			t.Errorf("expected user's active vault to be empty string after deleting the vault, got %s", newUserData.ActiveVault)
		}
	})
}
