// test/user_repository_helpers.go
package test

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// SetupUserRepositoryTest создает мок базы данных и репозиторий для тестирования
func SetupUserRepositoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := postgres.NewUserRepository(db)
	return db, mock, repo
}

// TestUserGetHelper тестирует метод Get
func TestUserGetHelper(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	userID := uuid.New()
	registrationDate := time.Now().Add(-30 * 24 * time.Hour)
	lastActivity := time.Now().Add(-2 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "age", "gender", "city", "registration_date", "last_activity"}).
		AddRow(userID, "user@example.com", "+1234567890", 30, "male", "New York", registrationDate, lastActivity)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.Get(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "user@example.com", user.Email)
	assert.Equal(t, "+1234567890", user.Phone)
	assert.Equal(t, 30, user.Age)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, "New York", user.City)
	assert.Equal(t, registrationDate.Unix(), user.RegistrationDate.Unix())
	assert.Equal(t, lastActivity.Unix(), user.LastActivity.Unix())
}

// TestUserCreateHelper тестирует метод Create
func TestUserCreateHelper(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	userID := uuid.New()
	registrationDate := time.Now()
	lastActivity := time.Now()

	user := &entities.User{
		Email:  "newuser@example.com",
		Phone:  "+9876543210",
		Age:    25,
		Gender: "female",
		City:   "San Francisco",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Email, user.Phone, user.Age, user.Gender, user.City).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "last_activity"}).
			AddRow(userID, registrationDate, lastActivity))

	err := repo.Create(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, registrationDate.Unix(), user.RegistrationDate.Unix())
	assert.Equal(t, lastActivity.Unix(), user.LastActivity.Unix())
}

// TestUserUpdateHelper тестирует метод Update
func TestUserUpdateHelper(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	user := &entities.User{
		ID:     uuid.New(),
		Email:  "updated@example.com",
		Phone:  "+1122334455",
		Age:    35,
		Gender: "male",
		City:   "Chicago",
	}

	mock.ExpectExec("UPDATE users SET").
		WithArgs(user.Email, user.Phone, user.Age, user.Gender, user.City, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(ctx, user)

	assert.NoError(t, err)
}

// TestUserDeleteHelper тестирует метод Delete
func TestUserDeleteHelper(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	userID := uuid.New()

	mock.ExpectExec("DELETE FROM users WHERE id = (.+)").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(ctx, userID)

	assert.NoError(t, err)
}

// TestUserListHelper тестирует метод List
func TestUserListHelper(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	limit := 10
	offset := 0

	user1ID := uuid.New()
	user2ID := uuid.New()
	registrationDate1 := time.Now().Add(-60 * 24 * time.Hour)
	registrationDate2 := time.Now().Add(-30 * 24 * time.Hour)
	lastActivity1 := time.Now().Add(-5 * 24 * time.Hour)
	lastActivity2 := time.Now().Add(-2 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "age", "gender", "city", "registration_date", "last_activity"}).
		AddRow(user1ID, "user1@example.com", "+1111111111", 28, "female", "Boston", registrationDate1, lastActivity1).
		AddRow(user2ID, "user2@example.com", "+2222222222", 42, "male", "Miami", registrationDate2, lastActivity2)

	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY registration_date DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	users, err := repo.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// First user
	assert.Equal(t, user1ID, users[0].ID)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.Equal(t, "+1111111111", users[0].Phone)
	assert.Equal(t, 28, users[0].Age)
	assert.Equal(t, "female", users[0].Gender)
	assert.Equal(t, "Boston", users[0].City)

	// Second user
	assert.Equal(t, user2ID, users[1].ID)
	assert.Equal(t, "user2@example.com", users[1].Email)
	assert.Equal(t, "+2222222222", users[1].Phone)
	assert.Equal(t, 42, users[1].Age)
	assert.Equal(t, "male", users[1].Gender)
	assert.Equal(t, "Miami", users[1].City)
}
