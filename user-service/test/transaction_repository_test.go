// test/transaction_repository_test.go
package test

import (
	"testing"
)

func TestTransactionRepository_GetByPeriod_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupTransactionRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestGetByPeriodHelper(t, repo, mock)
}

func TestTransactionRepository_GetByUserID_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupTransactionRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestGetByUserIDHelper(t, repo, mock)
}

func TestTransactionRepository_Create_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupTransactionRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestCreateHelper(t, repo, mock)
}
