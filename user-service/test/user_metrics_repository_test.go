// test/user_metrics_repository_test.go
package test

import (
	"testing"
)

func TestUserMetricsRepository_Get_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserMetricsRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserMetricsGetHelper(t, repo, mock)
}

func TestUserMetricsRepository_Create_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserMetricsRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserMetricsCreateHelper(t, repo, mock)
}

func TestUserMetricsRepository_Update_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserMetricsRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserMetricsUpdateHelper(t, repo, mock)
}

func TestUserMetricsRepository_CalculateMetrics_Standalone(t *testing.T) {
	// Setup
	db, mock, repo := SetupUserMetricsRepositoryTest(t)
	defer db.Close()

	// Execute test using helper
	TestUserMetricsCalculateMetricsHelper(t, repo, mock)
}
