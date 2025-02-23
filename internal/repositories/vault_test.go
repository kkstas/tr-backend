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

func TestVaultRepo_CreateOne(t *testing.T) {
	t.Parallel()

	t.Run("creates new vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)
		vaultRepo := repositories.NewVaultRepo(db)

		_, err := vaultRepo.CreateOne(ctx, user.ID, models.VaultRoleOwner, "some name")
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultRepo.FindAll(ctx, user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)
	})
}

func TestVaultRepo_FindAll(t *testing.T) {
	t.Parallel()

	t.Run("finds all vaults for user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)

		vaultRepo := repositories.NewVaultRepo(db)

		userID := user.ID
		userRole := models.VaultRoleOwner
		vaultName := "some name"

		_, err := vaultRepo.CreateOne(ctx, userID, userRole, vaultName)
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultRepo.FindAll(ctx, userID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 1)
		testutils.AssertEqual(t, foundVaults[0].UserRole, userRole)
		testutils.AssertEqual(t, foundVaults[0].Name, vaultName)
	})

	t.Run("returns empty array if no vaults are found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)

		vaultRepo := repositories.NewVaultRepo(db)

		foundVaults, err := vaultRepo.FindAll(ctx, uuid.New().String())
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 0)
	})
}

func TestVaultRepo_FindOneByID(t *testing.T) {
	t.Parallel()

	t.Run("finds one vault by ID", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)

		vaultRepo := repositories.NewVaultRepo(db)

		userID := user.ID
		userRole := models.VaultRoleOwner
		vaultName := "some name"

		vaultID, err := vaultRepo.CreateOne(ctx, userID, userRole, vaultName)
		testutils.AssertNoError(t, err)

		vault, err := vaultRepo.FindOneByID(ctx, userID, vaultID)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, vault.ID, vaultID)
		testutils.AssertEqual(t, vault.Name, vaultName)
		testutils.AssertEqual(t, vault.UserRole, userRole)
	})
}

func TestVaultRepo_FindOneByName(t *testing.T) {
	t.Parallel()

	t.Run("finds one vault by name", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)

		vaultRepo := repositories.NewVaultRepo(db)

		userID := user.ID
		userRole := models.VaultRoleOwner
		vaultName := "some name"

		vaultID, err := vaultRepo.CreateOne(ctx, userID, userRole, vaultName)
		testutils.AssertNoError(t, err)

		vault, err := vaultRepo.FindOneByName(ctx, userID, vaultName)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, vault.ID, vaultID)
		testutils.AssertEqual(t, vault.Name, vaultName)
		testutils.AssertEqual(t, vault.UserRole, userRole)
	})

	t.Run("returns error when no vault is found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		vaultRepo := repositories.NewVaultRepo(testutils.OpenTestDB(t, ctx))

		_, err := vaultRepo.FindOneByName(ctx, uuid.New().String(), uuid.New().String())
		if err == nil {
			t.Error("expected an error but didn't get one")
		}

		want := repositories.ErrVaultNotFound
		if !errors.Is(err, want) {
			t.Errorf("expected error %q, got %v", want, err)
		}
	})
}

func TestVaultRepo_DeleteOneByID(t *testing.T) {
	t.Parallel()

	t.Run("deletes existing vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)

		vaultRepo := repositories.NewVaultRepo(db)

		vaultID, err := vaultRepo.CreateOne(ctx, user.ID, models.VaultRoleOwner, "some name")
		testutils.AssertNoError(t, err)

		err = vaultRepo.DeleteOneByID(ctx, vaultID)
		testutils.AssertNoError(t, err)

		foundVaults, err := vaultRepo.FindAll(ctx, user.ID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, len(foundVaults), 0)
	})
}

func TestVaultRepo_AddUser(t *testing.T) {
	t.Parallel()

	t.Run("adds user to vault", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		db := testutils.OpenTestDB(t, ctx)
		user := testutils.CreateTestUser(t, db)

		vaultRepo := repositories.NewVaultRepo(db)

		userID := user.ID
		userRole := models.VaultRoleOwner
		vaultName := "some name"

		vaultID, err := vaultRepo.CreateOne(ctx, userID, userRole, vaultName)
		testutils.AssertNoError(t, err)

		invitee := testutils.CreateTestUser(t, db)

		inviteeRole := models.VaultRoleEditor
		err = vaultRepo.AddUser(ctx, vaultID, invitee.ID, inviteeRole)
		if err != nil {
			t.Errorf("didn't expect an error but got one: %v", err)
		}

		inviteeVaultWithRole, err := vaultRepo.FindOneByID(ctx, invitee.ID, vaultID)
		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, inviteeVaultWithRole.Name, vaultName)
		testutils.AssertEqual(t, inviteeVaultWithRole.UserRole, inviteeRole)
	})
}
