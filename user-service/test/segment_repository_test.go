// test/segment_repository_test.go
package test

import (
	"testing"
)

func TestSegmentRepository_Get_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupSegmentRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestSegmentGetHelper(t, repo, mock)
}

func TestSegmentRepository_Create_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupSegmentRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestSegmentCreateHelper(t, repo, mock)
}

func TestSegmentRepository_Update_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupSegmentRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestSegmentUpdateHelper(t, repo, mock)
}

func TestSegmentRepository_GetByType_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupSegmentRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestSegmentGetByTypeHelper(t, repo, mock)
}
