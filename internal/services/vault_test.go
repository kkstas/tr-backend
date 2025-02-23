package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kkstas/tr-backend/internal/models"
	"github.com/kkstas/tr-backend/internal/services"
	"github.com/kkstas/tr-backend/internal/testutils"
)

func TestVaultService_FindOneByID(t *testing.T) {
	t.Parallel()

	t.Run("finds one vault for user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultService := testutils.NewTestVaultService(db)
		createdUser := testutils.CreateTestUser(t, db)

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

	t.Run("returns error when no vault with given id is found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultService := testutils.NewTestVaultService(db)
		createdUser := testutils.CreateTestUser(t, db)

		_, err := vaultService.FindOneByID(ctx, createdUser.ID, uuid.New().String())
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := services.ErrVaultNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
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
		createdUser := testutils.CreateTestUser(t, db)

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

	t.Run("returns error if vault with given name already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultService := testutils.NewTestVaultService(db)
		createdUser := testutils.CreateTestUser(t, db)

		vaultName := "some vault"

		err := vaultService.CreateOne(ctx, createdUser.ID, vaultName)
		testutils.AssertNoError(t, err)

		err = vaultService.CreateOne(ctx, createdUser.ID, vaultName)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := services.ErrVaultWithThatNameAlreadyExists
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestVaultService_DeleteOneByID(t *testing.T) {
	t.Parallel()

	t.Run("deletes vault and clears user's active_vault property if deleted vault was the active one", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		createdUser := testutils.CreateTestUser(t, db)
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

func TestVaultService_AddUser(t *testing.T) {
	t.Parallel()

	t.Run("adds user to vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultOwner := testutils.CreateTestUser(t, db)
		invitee := testutils.CreateTestUser(t, db)
		vaultService := testutils.NewTestVaultService(db)

		err := vaultService.CreateOne(ctx, vaultOwner.ID, "vault name")
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultService.FindAll(ctx, vaultOwner.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)

		err = vaultService.AddUser(ctx, vaultOwner.ID, invitee.ID, foundVaults[0].ID, models.VaultRoleEditor)
		if err != nil {
			t.Errorf("didn't expect an error, but got one: %v", err)
		}

		inviteeVault, err := vaultService.FindOneByID(ctx, invitee.ID, foundVaults[0].ID)
		testutils.AssertNoError(t, err)
		testutils.AssertNotEmpty(t, inviteeVault)
		testutils.AssertEqual(t, inviteeVault.UserRole, models.VaultRoleEditor)
	})

	t.Run("returns error if inviter is not a vault owner", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultOwner := testutils.CreateTestUser(t, db)
		inviter := testutils.CreateTestUser(t, db)
		invitee := testutils.CreateTestUser(t, db)
		vaultService := testutils.NewTestVaultService(db)

		err := vaultService.CreateOne(ctx, vaultOwner.ID, "vault name")
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultService.FindAll(ctx, vaultOwner.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)

		err = vaultService.AddUser(ctx, vaultOwner.ID, inviter.ID, foundVaults[0].ID, models.VaultRoleEditor)
		testutils.AssertNoError(t, err)

		err = vaultService.AddUser(ctx, inviter.ID, invitee.ID, foundVaults[0].ID, models.VaultRoleEditor)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := services.ErrInsufficientVaultPermissions
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})

	t.Run("returns error if inviter does not belong to this vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		vaultOwner := testutils.CreateTestUser(t, db)
		inviter := testutils.CreateTestUser(t, db)
		invitee := testutils.CreateTestUser(t, db)
		vaultService := testutils.NewTestVaultService(db)

		err := vaultService.CreateOne(ctx, vaultOwner.ID, "vault name")
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultService.FindAll(ctx, vaultOwner.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)

		err = vaultService.AddUser(ctx, inviter.ID, invitee.ID, foundVaults[0].ID, models.VaultRoleEditor)
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := services.ErrVaultNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}
