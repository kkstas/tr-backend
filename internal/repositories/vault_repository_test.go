package repositories_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateAndFindAllVaults(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := testutils.OpenTestDB(t, ctx)
	vaultRepo := repositories.NewVaultRepo(db)

	name := "vault"

	err := vaultRepo.CreateOne(context.Background(), name)
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	foundVaults, err := vaultRepo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	testutils.AssertEqual(t, len(foundVaults), 1)
	testutils.AssertEqual(t, foundVaults[0].Name, name)
	testutils.AssertValidDate(t, foundVaults[0].CreatedAt)
	if err := uuid.Validate(foundVaults[0].ID); err != nil {
		t.Errorf("expected id to be valid uuid, got error: %v", err)
	}
}
