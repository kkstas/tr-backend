package repositories_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateAndFindAllUsers(t *testing.T) {
	db, cleanup := testutils.OpenTestDB(t)
	defer cleanup()
	userRepo := repositories.NewUserRepo(db)

	firstName := "John"
	lastName := "Doe"
	email := "john.doe@email.com"

	err := userRepo.CreateOne(firstName, lastName, email)
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	foundUsers, err := userRepo.FindAll()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	testutils.AssertEqual(t, len(foundUsers), 1)
	testutils.AssertEqual(t, foundUsers[0].FirstName, firstName)
	testutils.AssertEqual(t, foundUsers[0].LastName, lastName)
	testutils.AssertEqual(t, foundUsers[0].Email, email)
	testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
	if err := uuid.Validate(foundUsers[0].Id); err != nil {
		t.Errorf("expected id to be valid uuid, got error: %v", err)
	}
}
