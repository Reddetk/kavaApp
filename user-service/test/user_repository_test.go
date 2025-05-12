// test/user_repository_test.go
package test

import (
	"testing"
)

func TestUserRepository_Get_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserGetHelper(t, repo, mock)
}

func TestUserRepository_Create_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserCreateHelper(t, repo, mock)
}

func TestUserRepository_Update_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserUpdateHelper(t, repo, mock)
}

func TestUserRepository_Delete_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserDeleteHelper(t, repo, mock)
}

func TestUserRepository_List_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserListHelper(t, repo, mock)
}
